package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/luckermt/shared/logger"
	"go.uber.org/zap"
)

// JWTClaims структура для хранения данных в токене
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateJWTToken создает новый JWT токен
func GenerateJWTToken(userID, username, secret string) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		logger.Log.Error("Failed to sign JWT token", zap.Error(err))
		return "", err
	}

	return signedToken, nil
}

// ParseJWTToken проверяет и парсит JWT токен
func ParseJWTToken(tokenString, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		logger.Log.Error("Failed to parse JWT token", zap.Error(err))
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}
