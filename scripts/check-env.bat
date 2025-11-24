@echo off
REM Environment Checker for Windows
REM Validates your setup before running

echo ==========================================
echo   Environment Check
echo ==========================================
echo.

set ERRORS=0
set WARNINGS=0

REM Check Go installation
echo Checking Go installation...
where go >nul 2>nul
if %ERRORLEVEL% EQU 0 (
    for /f "tokens=3" %%i in ('go version') do (
        echo [OK] Found Go %%i
        goto :node_check
    )
) else (
    echo [FAIL] Go not found
    echo   Please install Go 1.21+ from https://go.dev/dl/
    set /a ERRORS+=1
)

:node_check
REM Check Node.js installation
echo Checking Node.js installation...
where node >nul 2>nul
if %ERRORLEVEL% EQU 0 (
    for /f "tokens=*" %%i in ('node --version') do (
        echo [OK] Found Node.js %%i
    )
) else (
    echo [FAIL] Node.js not found
    echo   Please install Node.js 18+ from https://nodejs.org/
    set /a ERRORS+=1
)

REM Check npm installation
echo Checking npm installation...
where npm >nul 2>nul
if %ERRORLEVEL% EQU 0 (
    for /f "tokens=*" %%i in ('npm --version') do (
        echo [OK] Found npm v%%i
    )
) else (
    echo [FAIL] npm not found
    set /a ERRORS+=1
)

REM Check .env file
echo Checking .env file...
if exist .env (
    echo [OK] Found .env

    REM Check Google OAuth
    findstr /C:"GOOGLE_CLIENT_ID=" .env | findstr /V "your-google-client-id" >nul
    if %ERRORLEVEL% EQU 0 (
        echo   [OK] Google OAuth configured
    ) else (
        echo   [WARN] Google OAuth not configured
        set /a WARNINGS+=1
    )

    REM Check Microsoft OAuth
    findstr /C:"MICROSOFT_CLIENT_ID=" .env | findstr /V "your-microsoft-client-id" >nul
    if %ERRORLEVEL% EQU 0 (
        echo   [OK] Microsoft OAuth configured
    ) else (
        echo   [WARN] Microsoft OAuth not configured
        set /a WARNINGS+=1
    )

    REM Check OpenAI
    findstr /C:"OPENAI_API_KEY=" .env | findstr /V "your-openai-api-key" >nul
    if %ERRORLEVEL% EQU 0 (
        echo   [OK] OpenAI API configured
    ) else (
        echo   [WARN] OpenAI API not configured (AI features won't work^)
        set /a WARNINGS+=1
    )
) else (
    echo [FAIL] .env file not found
    echo   Run: copy .env.example .env
    echo   Then configure your credentials
    set /a ERRORS+=1
)

echo.
echo ==========================================
echo   Results
echo ==========================================
echo.

if %ERRORS% EQU 0 (
    if %WARNINGS% EQU 0 (
        echo [OK] All checks passed! Ready to run.
        echo.
        echo Start the application with:
        echo   scripts\quickstart.bat
        echo   or
        echo   make dev
        exit /b 0
    ) else (
        echo [WARN] Found %WARNINGS% warning(s^)
        echo.
        echo You can still run the application, but some features may not work.
        echo.
        echo To fix warnings, run:
        echo   scripts\setup-oauth.bat
        exit /b 0
    )
) else (
    echo [FAIL] Found %ERRORS% error(s^) and %WARNINGS% warning(s^)
    echo.
    echo Please fix the errors above before running.
    pause
    exit /b 1
)
