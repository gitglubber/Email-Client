@echo off
REM Windows Batch wrapper for OAuth Setup
REM This calls the PowerShell script

echo ==========================================
echo   AI Email Client - OAuth Setup Helper
echo ==========================================
echo.

REM Check if PowerShell is available
where powershell >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Error: PowerShell is not available
    echo Please install PowerShell or run the script manually
    pause
    exit /b 1
)

REM Run the PowerShell script
powershell -ExecutionPolicy Bypass -File "%~dp0setup-oauth.ps1"

pause
