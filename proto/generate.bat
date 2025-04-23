@echo off
SETLOCAL

:: Проверка наличия protoc
where protoc >nul 2>&1
if %errorlevel% neq 0 (
    echo Error: protoc not found in PATH
    echo Install protoc from https://github.com/protocolbuffers/protobuf/releases
    pause
    exit /b 1
)

:: Переход в директорию скрипта
pushd %~dp0

:: Генерация auth.proto
if exist "auth\auth.proto" (
    protoc --go_out=. --go_opt=paths=source_relative ^
           --go-grpc_out=. --go-grpc_opt=paths=source_relative ^
           auth\auth.proto
) else (
    echo Error: auth\auth.proto not found
)

:: Генерация forum.proto
if exist "forum\forum.proto" (
    protoc --go_out=. --go_opt=paths=source_relative ^
           --go-grpc_out=. --go-grpc_opt=paths=source_relative ^
           forum\forum.proto
) else (
    echo Error: forum\forum.proto not found
)

:: Возврат из директории
popd

echo Generation completed
pause