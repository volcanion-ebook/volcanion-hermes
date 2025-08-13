# Volcanion Hermes Setup Script for Windows
# Run this script in PowerShell as Administrator

Write-Host "Setting up Volcanion Hermes development environment..." -ForegroundColor Green

# Check if Go is installed
$goVersion = $null
try {
    $goVersion = go version 2>$null
    Write-Host "Go is already installed: $goVersion" -ForegroundColor Green
} catch {
    Write-Host "Go is not installed. Please install Go from https://golang.org/dl/" -ForegroundColor Red
    Write-Host "1. Download Go installer for Windows"
    Write-Host "2. Run the installer"
    Write-Host "3. Restart PowerShell"
    Write-Host "4. Run this script again"
    exit 1
}

# Check if Docker is installed
$dockerVersion = $null
try {
    $dockerVersion = docker --version 2>$null
    Write-Host "Docker is already installed: $dockerVersion" -ForegroundColor Green
} catch {
    Write-Host "Docker is not installed. Please install Docker Desktop from https://www.docker.com/products/docker-desktop" -ForegroundColor Yellow
}

# Navigate to project directory
$projectDir = "G:\Github\volcanion-ebook\volcanion-hermes"
Set-Location $projectDir

Write-Host "Current directory: $(Get-Location)" -ForegroundColor Blue

# Initialize Go module if go.mod doesn't exist
if (-Not (Test-Path "go.mod")) {
    Write-Host "Initializing Go module..." -ForegroundColor Yellow
    go mod init github.com/volcanion/volcanion-hermes
}

# Download dependencies
Write-Host "Downloading Go dependencies..." -ForegroundColor Yellow
go mod tidy

# Create required directories
$dirs = @("logs", "uploads", "tmp")
foreach ($dir in $dirs) {
    if (-Not (Test-Path $dir)) {
        New-Item -ItemType Directory -Path $dir
        Write-Host "Created directory: $dir" -ForegroundColor Green
    }
}

# Start development services with Docker
Write-Host "Starting development services..." -ForegroundColor Yellow

# Start MongoDB
Write-Host "Starting MongoDB..." -ForegroundColor Blue
docker run --name volcanion-mongodb -d -p 27017:27017 mongo:6.0 2>$null
if ($LASTEXITCODE -ne 0) {
    Write-Host "MongoDB container already exists or failed to start" -ForegroundColor Yellow
    docker start volcanion-mongodb 2>$null
}

# Start MinIO
Write-Host "Starting MinIO..." -ForegroundColor Blue
docker run --name volcanion-minio -d -p 9000:9000 -p 9001:9001 -e "MINIO_ROOT_USER=minioadmin" -e "MINIO_ROOT_PASSWORD=minioadmin" minio/minio server /data --console-address ":9001" 2>$null
if ($LASTEXITCODE -ne 0) {
    Write-Host "MinIO container already exists or failed to start" -ForegroundColor Yellow
    docker start volcanion-minio 2>$null
}

# Wait for services to start
Write-Host "Waiting for services to start..." -ForegroundColor Yellow
Start-Sleep 10

Write-Host "Setup completed!" -ForegroundColor Green
Write-Host ""
Write-Host "Services:" -ForegroundColor Cyan
Write-Host "- MongoDB: mongodb://localhost:27017" -ForegroundColor White
Write-Host "- MinIO Console: http://localhost:9001 (admin/minioadmin)" -ForegroundColor White
Write-Host "- MinIO API: http://localhost:9000" -ForegroundColor White
Write-Host ""
Write-Host "To start the application:" -ForegroundColor Cyan
Write-Host "go run cmd/server/main.go" -ForegroundColor White
Write-Host ""
Write-Host "Or use Make commands:" -ForegroundColor Cyan
Write-Host "make dev           # Run in development mode" -ForegroundColor White
Write-Host "make build         # Build the application" -ForegroundColor White
Write-Host "make test          # Run tests" -ForegroundColor White
Write-Host "make dev-up        # Start all services" -ForegroundColor White
Write-Host "make dev-down      # Stop all services" -ForegroundColor White
