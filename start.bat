@echo off
echo Starting Volcanion Hermes Ebook Management System...
echo.

REM Check if Go is installed
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo Error: Go is not installed or not in PATH
    echo Please install Go from https://golang.org/dl/
    pause
    exit /b 1
)

REM Navigate to project directory
cd /d "G:\Github\volcanion-ebook\volcanion-hermes"

REM Check if .env exists
if not exist ".env" (
    echo Error: .env file not found
    echo Please copy and configure the .env file
    pause
    exit /b 1
)

REM Start MongoDB if not running
echo Checking MongoDB...
docker ps | findstr volcanion-mongodb >nul 2>&1
if %errorlevel% neq 0 (
    echo Starting MongoDB...
    docker run --name volcanion-mongodb -d -p 27017:27017 mongo:6.0 >nul 2>&1
    if %errorlevel% neq 0 (
        docker start volcanion-mongodb >nul 2>&1
    )
)

REM Start MinIO if not running
echo Checking MinIO...
docker ps | findstr volcanion-minio >nul 2>&1
if %errorlevel% neq 0 (
    echo Starting MinIO...
    docker run --name volcanion-minio -d -p 9000:9000 -p 9001:9001 -e "MINIO_ROOT_USER=minioadmin" -e "MINIO_ROOT_PASSWORD=minioadmin" minio/minio server /data --console-address ":9001" >nul 2>&1
    if %errorlevel% neq 0 (
        docker start volcanion-minio >nul 2>&1
    )
)

REM Wait for services
echo Waiting for services to start...
timeout /t 5 /nobreak >nul

REM Download dependencies if needed
if not exist "go.sum" (
    echo Downloading dependencies...
    go mod tidy
)

echo.
echo Starting the application...
echo Server will be available at: http://localhost:8080
echo MongoDB: mongodb://localhost:27017
echo MinIO Console: http://localhost:9001 (admin/minioadmin)
echo.

REM Start the application
go run cmd/server/main.go

pause
