@echo off
REM Quick Start Script for Windows
REM Sets up and runs the application

echo ==========================================
echo   AI Email Client - Quick Start
echo ==========================================
echo.

REM Check if .env exists
if not exist .env (
    echo No .env file found. Running OAuth setup...
    echo.
    call "%~dp0setup-oauth.bat"
    if %ERRORLEVEL% NEQ 0 (
        echo Setup failed. Exiting...
        pause
        exit /b 1
    )
) else (
    REM Check if credentials are set
    findstr /C:"your-google-client-id" .env >nul
    if %ERRORLEVEL% EQU 0 (
        echo OAuth credentials not configured.
        echo.
        set /p answer="Would you like to run the setup wizard? (y/n): "
        if /i "%answer%"=="y" (
            call "%~dp0setup-oauth.bat"
        ) else (
            echo Please configure .env manually or run: scripts\setup-oauth.bat
            pause
            exit /b 1
        )
    )
)

REM Check if Go is installed
where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Error: Go is not installed. Please install Go 1.21+ from https://go.dev/dl/
    pause
    exit /b 1
)

REM Check if Node.js is installed
where node >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Error: Node.js is not installed. Please install Node.js 18+ from https://nodejs.org/
    pause
    exit /b 1
)

echo.
echo Installing dependencies...
echo.

REM Install Go dependencies
echo Installing Go dependencies...
go mod download
if %ERRORLEVEL% NEQ 0 (
    echo Failed to download Go dependencies
    pause
    exit /b 1
)

REM Install frontend dependencies
echo Installing frontend dependencies...
cd web
call npm install --silent
if %ERRORLEVEL% NEQ 0 (
    echo Failed to install npm dependencies
    cd ..
    pause
    exit /b 1
)
cd ..

echo.
echo ==========================================
echo   Starting Application
echo ==========================================
echo.
echo Backend: http://localhost:8080
echo Frontend: http://localhost:3000
echo.
echo Press Ctrl+C to stop both servers
echo.

REM Start both servers (simplified version)
echo Starting servers...
echo.
echo NOTE: For full development mode with hot reload,
echo       run in separate terminals:
echo       Terminal 1: go run cmd\server\main.go
echo       Terminal 2: cd web ^&^& npm start
echo.
echo Starting backend...
start "Email Client Backend" cmd /c "go run cmd\server\main.go"
timeout /t 3 /nobreak >nul
echo Starting frontend...
cd web
npm start
