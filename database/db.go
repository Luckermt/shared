package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/luckermt/shared/config"
	"github.com/luckermt/shared/logger"
	"go.uber.org/zap"
)

// DB представляет обёртку вокруг sql.DB с дополнительными методами
type DB struct {
	*sql.DB
}

// NewPostgresConnection создает новое подключение к PostgreSQL
func NewPostgresConnection(cfg *config.Config) (*DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Устанавливаем разумные ограничения на подключения
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Проверяем подключение
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Log.Info("Successfully connected to PostgreSQL database")

	return &DB{db}, nil
}

// WithTransaction выполняет переданную функцию в транзакции
func (db *DB) WithTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			logger.Log.Error("failed to rollback transaction",
				zap.Error(rbErr),
				zap.NamedError("original_error", err))
			return fmt.Errorf("transaction rollback error: %v (original error: %w)", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Close закрывает подключение к базе данных
func (db *DB) Close() error {
	if err := db.DB.Close(); err != nil {
		logger.Log.Error("failed to close database connection", zap.Error(err))
		return err
	}
	logger.Log.Info("Database connection closed")
	return nil
}
