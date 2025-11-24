# OAuth Setup Helper Script for Windows
# This script guides you through setting up OAuth credentials

$ErrorActionPreference = "Stop"

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "  AI Email Client - OAuth Setup Helper" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

# Check if .env exists
if (Test-Path ".env") {
    Write-Host "Warning: .env file already exists. This will update it." -ForegroundColor Yellow
    $continue = Read-Host "Continue? (y/n)"
    if ($continue -ne "y") {
        exit 1
    }
}

# Copy template if needed
if (-not (Test-Path ".env")) {
    Copy-Item ".env.example" ".env"
    Write-Host "Created .env file from template" -ForegroundColor Green
}

Write-Host ""
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "  Step 1: Google OAuth Setup" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Follow these steps:"
Write-Host "1. Go to https://console.cloud.google.com/"
Write-Host "2. Create a new project or select existing"
Write-Host "3. Enable Gmail API and Google Calendar API"
Write-Host "4. Go to 'Credentials' → 'Create Credentials' → 'OAuth 2.0 Client ID'"
Write-Host "5. Configure OAuth consent screen (add email, profile, gmail, calendar scopes)"
Write-Host "6. Create Web application credentials"
Write-Host "7. Add redirect URI: http://localhost:8080/auth/google/callback"
Write-Host ""
Write-Host "Press Enter when you have your Google credentials ready..." -ForegroundColor Blue
Read-Host

Write-Host ""
$GOOGLE_CLIENT_ID = Read-Host "Enter Google Client ID"
$GOOGLE_CLIENT_SECRET = Read-Host "Enter Google Client Secret"

# Update .env file
$envContent = Get-Content ".env"
$envContent = $envContent -replace "^GOOGLE_CLIENT_ID=.*", "GOOGLE_CLIENT_ID=$GOOGLE_CLIENT_ID"
$envContent = $envContent -replace "^GOOGLE_CLIENT_SECRET=.*", "GOOGLE_CLIENT_SECRET=$GOOGLE_CLIENT_SECRET"
$envContent | Set-Content ".env"

Write-Host "✓ Google OAuth configured" -ForegroundColor Green

Write-Host ""
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "  Step 2: Microsoft OAuth Setup" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Follow these steps:"
Write-Host "1. Go to https://portal.azure.com/"
Write-Host "2. Navigate to 'Azure Active Directory' → 'App registrations'"
Write-Host "3. Click 'New registration'"
Write-Host "4. Name: AI Email Client"
Write-Host "5. Supported account types: Personal Microsoft accounts"
Write-Host "6. Redirect URI: http://localhost:8080/auth/microsoft/callback"
Write-Host "7. After creation, go to 'Certificates & secrets' → Create new client secret"
Write-Host "8. Go to 'API permissions' → Add: User.Read, Mail.*, Calendars.ReadWrite"
Write-Host ""
Write-Host "Press Enter when you have your Microsoft credentials ready..." -ForegroundColor Blue
Read-Host

Write-Host ""
$MICROSOFT_CLIENT_ID = Read-Host "Enter Microsoft Client ID"
$MICROSOFT_CLIENT_SECRET = Read-Host "Enter Microsoft Client Secret"

# Update .env file
$envContent = Get-Content ".env"
$envContent = $envContent -replace "^MICROSOFT_CLIENT_ID=.*", "MICROSOFT_CLIENT_ID=$MICROSOFT_CLIENT_ID"
$envContent = $envContent -replace "^MICROSOFT_CLIENT_SECRET=.*", "MICROSOFT_CLIENT_SECRET=$MICROSOFT_CLIENT_SECRET"
$envContent | Set-Content ".env"

Write-Host "✓ Microsoft OAuth configured" -ForegroundColor Green

Write-Host ""
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "  Step 3: OpenAI Configuration" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Choose your AI provider:"
Write-Host "1) OpenAI (GPT-4)"
Write-Host "2) Azure OpenAI"
Write-Host "3) Local LLM (LM Studio/Ollama)"
Write-Host "4) Other OpenAI-compatible endpoint"
Write-Host ""
$AI_CHOICE = Read-Host "Select option (1-4)"

switch ($AI_CHOICE) {
    "1" {
        Write-Host ""
        Write-Host "Get your API key from: https://platform.openai.com/api-keys"
        $OPENAI_API_KEY = Read-Host "Enter OpenAI API Key"
        $OPENAI_BASE_URL = "https://api.openai.com/v1"
        $OPENAI_MODEL = "gpt-4-turbo-preview"
    }
    "2" {
        Write-Host ""
        $OPENAI_API_KEY = Read-Host "Enter Azure OpenAI API Key"
        $AZURE_RESOURCE = Read-Host "Enter Azure Resource Name"
        $OPENAI_BASE_URL = "https://$AZURE_RESOURCE.openai.azure.com/"
        $OPENAI_MODEL = Read-Host "Enter deployment name"
    }
    "3" {
        Write-Host ""
        Write-Host "Make sure LM Studio or Ollama is running!"
        Write-Host "LM Studio: Start server in Settings"
        Write-Host "Ollama: Run 'ollama serve'"
        $LOCAL_URL = Read-Host "Enter local server URL (default: http://localhost:1234/v1)"
        if ([string]::IsNullOrWhiteSpace($LOCAL_URL)) {
            $LOCAL_URL = "http://localhost:1234/v1"
        }
        $OPENAI_API_KEY = "not-needed"
        $OPENAI_BASE_URL = $LOCAL_URL
        $OPENAI_MODEL = Read-Host "Enter model name"
    }
    "4" {
        $OPENAI_API_KEY = Read-Host "Enter API Key"
        $OPENAI_BASE_URL = Read-Host "Enter Base URL"
        $OPENAI_MODEL = Read-Host "Enter Model name"
    }
}

# Update .env file
$envContent = Get-Content ".env"
$envContent = $envContent -replace "^OPENAI_API_KEY=.*", "OPENAI_API_KEY=$OPENAI_API_KEY"
$envContent = $envContent -replace "^OPENAI_BASE_URL=.*", "OPENAI_BASE_URL=$OPENAI_BASE_URL"
$envContent = $envContent -replace "^OPENAI_MODEL=.*", "OPENAI_MODEL=$OPENAI_MODEL"
$envContent | Set-Content ".env"

# Generate random session secret
$bytes = New-Object byte[] 32
$rng = [System.Security.Cryptography.RNGCryptoServiceProvider]::Create()
$rng.GetBytes($bytes)
$SESSION_SECRET = [Convert]::ToBase64String($bytes)

$envContent = Get-Content ".env"
$envContent = $envContent -replace "^SESSION_SECRET=.*", "SESSION_SECRET=$SESSION_SECRET"
$envContent | Set-Content ".env"

Write-Host ""
Write-Host "==========================================" -ForegroundColor Green
Write-Host "  ✓ Setup Complete!" -ForegroundColor Green
Write-Host "==========================================" -ForegroundColor Green
Write-Host ""
Write-Host "Your .env file has been configured with:"
Write-Host "  - Google OAuth credentials"
Write-Host "  - Microsoft OAuth credentials"
Write-Host "  - AI provider settings"
Write-Host "  - Secure session secret"
Write-Host ""
Write-Host "Next steps:"
Write-Host "  1. Install dependencies: make install-deps"
Write-Host "  2. Start the app: make dev"
Write-Host "  3. Open http://localhost:3000"
Write-Host ""
Write-Host "Happy emailing! 🚀" -ForegroundColor Blue
