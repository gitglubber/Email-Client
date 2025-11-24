# Environment Checker - Validates your setup before running
# PowerShell version

$ErrorActionPreference = "Stop"

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "  Environment Check" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

$ERRORS = 0
$WARNINGS = 0

# Check Go installation
Write-Host "Checking Go installation... " -NoNewline
try {
    $goVersion = go version 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ Found $goVersion" -ForegroundColor Green
    } else {
        throw
    }
} catch {
    Write-Host "✗ Go not found" -ForegroundColor Red
    Write-Host "  Please install Go 1.21+ from https://go.dev/dl/" -ForegroundColor Yellow
    $ERRORS++
}

# Check Node.js installation
Write-Host "Checking Node.js installation... " -NoNewline
try {
    $nodeVersion = node --version 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ Found $nodeVersion" -ForegroundColor Green
    } else {
        throw
    }
} catch {
    Write-Host "✗ Node.js not found" -ForegroundColor Red
    Write-Host "  Please install Node.js 18+ from https://nodejs.org/" -ForegroundColor Yellow
    $ERRORS++
}

# Check npm installation
Write-Host "Checking npm installation... " -NoNewline
try {
    $npmVersion = npm --version 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ Found v$npmVersion" -ForegroundColor Green
    } else {
        throw
    }
} catch {
    Write-Host "✗ npm not found" -ForegroundColor Red
    $ERRORS++
}

# Check .env file
Write-Host "Checking .env file... " -NoNewline
if (Test-Path ".env") {
    Write-Host "✓ Found" -ForegroundColor Green

    $envContent = Get-Content ".env" -Raw

    # Check Google OAuth
    Write-Host "  Checking Google OAuth... " -NoNewline
    if ($envContent -match "GOOGLE_CLIENT_ID=(.+)" -and $Matches[1] -ne "your-google-client-id" -and $Matches[1] -ne "") {
        Write-Host "✓ Configured" -ForegroundColor Green
    } else {
        Write-Host "✗ Not configured" -ForegroundColor Red
        $WARNINGS++
    }

    # Check Microsoft OAuth
    Write-Host "  Checking Microsoft OAuth... " -NoNewline
    if ($envContent -match "MICROSOFT_CLIENT_ID=(.+)" -and $Matches[1] -ne "your-microsoft-client-id" -and $Matches[1] -ne "") {
        Write-Host "✓ Configured" -ForegroundColor Green
    } else {
        Write-Host "✗ Not configured" -ForegroundColor Red
        $WARNINGS++
    }

    # Check OpenAI
    Write-Host "  Checking OpenAI API... " -NoNewline
    if ($envContent -match "OPENAI_API_KEY=(.+)" -and $Matches[1] -ne "your-openai-api-key" -and $Matches[1] -ne "") {
        Write-Host "✓ Configured" -ForegroundColor Green
    } else {
        Write-Host "⚠ Not configured (AI features won't work)" -ForegroundColor Yellow
        $WARNINGS++
    }

} else {
    Write-Host "✗ Not found" -ForegroundColor Red
    Write-Host "  Run: Copy-Item .env.example .env" -ForegroundColor Yellow
    Write-Host "  Then configure your credentials" -ForegroundColor Yellow
    $ERRORS++
}

# Check port availability
Write-Host "Checking port 8080 (backend)... " -NoNewline
$port8080 = Get-NetTCPConnection -LocalPort 8080 -ErrorAction SilentlyContinue
if ($port8080) {
    Write-Host "⚠ Port in use" -ForegroundColor Yellow
    $WARNINGS++
} else {
    Write-Host "✓ Available" -ForegroundColor Green
}

Write-Host "Checking port 3000 (frontend)... " -NoNewline
$port3000 = Get-NetTCPConnection -LocalPort 3000 -ErrorAction SilentlyContinue
if ($port3000) {
    Write-Host "⚠ Port in use" -ForegroundColor Yellow
    $WARNINGS++
} else {
    Write-Host "✓ Available" -ForegroundColor Green
}

Write-Host ""
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "  Results" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

if ($ERRORS -eq 0 -and $WARNINGS -eq 0) {
    Write-Host "✓ All checks passed! Ready to run." -ForegroundColor Green
    Write-Host ""
    Write-Host "Start the application with:"
    Write-Host "  .\scripts\quickstart.bat"
    Write-Host "  or"
    Write-Host "  make dev"
    exit 0
} elseif ($ERRORS -gt 0) {
    Write-Host "✗ Found $ERRORS error(s) and $WARNINGS warning(s)" -ForegroundColor Red
    Write-Host ""
    Write-Host "Please fix the errors above before running."
    exit 1
} else {
    Write-Host "⚠ Found $WARNINGS warning(s)" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "You can still run the application, but some features may not work."
    Write-Host ""
    Write-Host "To fix warnings, run:"
    Write-Host "  .\scripts\setup-oauth.bat"
    exit 0
}
