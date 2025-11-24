# Quick Start Script for Windows (PowerShell)
# Sets up and runs the application

$ErrorActionPreference = "Stop"

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "  AI Email Client - Quick Start" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

# Check if .env exists and has credentials
if (-not (Test-Path ".env")) {
    Write-Host "No .env file found. Running OAuth setup..." -ForegroundColor Yellow
    Write-Host ""
    & "$PSScriptRoot\setup-oauth.ps1"
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Setup failed. Exiting..." -ForegroundColor Red
        exit 1
    }
} else {
    # Check if credentials are set
    $envContent = Get-Content ".env" -Raw
    if ($envContent -match "your-google-client-id" -or $envContent -match "your-microsoft-client-id") {
        Write-Host "OAuth credentials not configured." -ForegroundColor Yellow
        Write-Host ""
        $answer = Read-Host "Would you like to run the setup wizard? (y/n)"
        if ($answer -eq "y") {
            & "$PSScriptRoot\setup-oauth.ps1"
        } else {
            Write-Host "Please configure .env manually or run: .\scripts\setup-oauth.ps1" -ForegroundColor Yellow
            exit 1
        }
    }
}

# Check if Go is installed
try {
    $null = go version 2>&1
    if ($LASTEXITCODE -ne 0) { throw }
} catch {
    Write-Host "Error: Go is not installed. Please install Go 1.21+ from https://go.dev/dl/" -ForegroundColor Red
    exit 1
}

# Check if Node.js is installed
try {
    $null = node --version 2>&1
    if ($LASTEXITCODE -ne 0) { throw }
} catch {
    Write-Host "Error: Node.js is not installed. Please install Node.js 18+ from https://nodejs.org/" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "Installing dependencies..." -ForegroundColor Cyan
Write-Host ""

# Install Go dependencies
Write-Host "Installing Go dependencies..." -ForegroundColor Yellow
go mod download
if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to download Go dependencies" -ForegroundColor Red
    exit 1
}

# Install frontend dependencies
Write-Host "Installing frontend dependencies..." -ForegroundColor Yellow
Push-Location web
npm install --silent
if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to install npm dependencies" -ForegroundColor Red
    Pop-Location
    exit 1
}
Pop-Location

Write-Host ""
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "  Starting Application" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Backend: http://localhost:8080" -ForegroundColor Green
Write-Host "Frontend: http://localhost:3000" -ForegroundColor Green
Write-Host ""
Write-Host "Press Ctrl+C to stop the servers" -ForegroundColor Yellow
Write-Host ""

# Start backend in background
Write-Host "Starting backend..." -ForegroundColor Cyan
$backend = Start-Process -FilePath "go" -ArgumentList "run", "cmd\server\main.go" -PassThru -NoNewWindow

# Wait a bit for backend to start
Start-Sleep -Seconds 3

# Start frontend
Write-Host "Starting frontend..." -ForegroundColor Cyan
Push-Location web
try {
    npm start
} finally {
    # Cleanup: Kill backend when frontend stops
    if (-not $backend.HasExited) {
        Write-Host "Stopping backend..." -ForegroundColor Yellow
        Stop-Process -Id $backend.Id -Force
    }
    Pop-Location
}
