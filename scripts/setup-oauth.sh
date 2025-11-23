#!/bin/bash

# OAuth Setup Helper Script
# This script guides you through setting up OAuth credentials

set -e

echo "=========================================="
echo "  AI Email Client - OAuth Setup Helper"
echo "=========================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if .env exists
if [ -f .env ]; then
    echo -e "${YELLOW}Warning: .env file already exists. This will update it.${NC}"
    read -p "Continue? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Copy template if needed
if [ ! -f .env ]; then
    cp .env.example .env
    echo -e "${GREEN}Created .env file from template${NC}"
fi

echo ""
echo "=========================================="
echo "  Step 1: Google OAuth Setup"
echo "=========================================="
echo ""
echo "Follow these steps:"
echo "1. Go to https://console.cloud.google.com/"
echo "2. Create a new project or select existing"
echo "3. Enable Gmail API and Google Calendar API"
echo "4. Go to 'Credentials' → 'Create Credentials' → 'OAuth 2.0 Client ID'"
echo "5. Configure OAuth consent screen (add email, profile, gmail, calendar scopes)"
echo "6. Create Web application credentials"
echo "7. Add redirect URI: http://localhost:8080/auth/google/callback"
echo ""
echo -e "${BLUE}Press Enter when you have your Google credentials ready...${NC}"
read

echo ""
read -p "Enter Google Client ID: " GOOGLE_CLIENT_ID
read -p "Enter Google Client Secret: " GOOGLE_CLIENT_SECRET

# Update .env file
sed -i.bak "s|^GOOGLE_CLIENT_ID=.*|GOOGLE_CLIENT_ID=$GOOGLE_CLIENT_ID|" .env
sed -i.bak "s|^GOOGLE_CLIENT_SECRET=.*|GOOGLE_CLIENT_SECRET=$GOOGLE_CLIENT_SECRET|" .env

echo -e "${GREEN}✓ Google OAuth configured${NC}"

echo ""
echo "=========================================="
echo "  Step 2: Microsoft OAuth Setup"
echo "=========================================="
echo ""
echo "Follow these steps:"
echo "1. Go to https://portal.azure.com/"
echo "2. Navigate to 'Azure Active Directory' → 'App registrations'"
echo "3. Click 'New registration'"
echo "4. Name: AI Email Client"
echo "5. Supported account types: Personal Microsoft accounts"
echo "6. Redirect URI: http://localhost:8080/auth/microsoft/callback"
echo "7. After creation, go to 'Certificates & secrets' → Create new client secret"
echo "8. Go to 'API permissions' → Add: User.Read, Mail.*, Calendars.ReadWrite"
echo ""
echo -e "${BLUE}Press Enter when you have your Microsoft credentials ready...${NC}"
read

echo ""
read -p "Enter Microsoft Client ID: " MICROSOFT_CLIENT_ID
read -p "Enter Microsoft Client Secret: " MICROSOFT_CLIENT_SECRET

# Update .env file
sed -i.bak "s|^MICROSOFT_CLIENT_ID=.*|MICROSOFT_CLIENT_ID=$MICROSOFT_CLIENT_ID|" .env
sed -i.bak "s|^MICROSOFT_CLIENT_SECRET=.*|MICROSOFT_CLIENT_SECRET=$MICROSOFT_CLIENT_SECRET|" .env

echo -e "${GREEN}✓ Microsoft OAuth configured${NC}"

echo ""
echo "=========================================="
echo "  Step 3: OpenAI Configuration"
echo "=========================================="
echo ""
echo "Choose your AI provider:"
echo "1) OpenAI (GPT-4)"
echo "2) Azure OpenAI"
echo "3) Local LLM (LM Studio/Ollama)"
echo "4) Other OpenAI-compatible endpoint"
echo ""
read -p "Select option (1-4): " AI_CHOICE

case $AI_CHOICE in
    1)
        echo ""
        echo "Get your API key from: https://platform.openai.com/api-keys"
        read -p "Enter OpenAI API Key: " OPENAI_API_KEY
        OPENAI_BASE_URL="https://api.openai.com/v1"
        OPENAI_MODEL="gpt-4-turbo-preview"
        ;;
    2)
        echo ""
        read -p "Enter Azure OpenAI API Key: " OPENAI_API_KEY
        read -p "Enter Azure Resource Name: " AZURE_RESOURCE
        OPENAI_BASE_URL="https://${AZURE_RESOURCE}.openai.azure.com/"
        read -p "Enter deployment name: " OPENAI_MODEL
        ;;
    3)
        echo ""
        echo "Make sure LM Studio or Ollama is running!"
        echo "LM Studio: Start server in Settings"
        echo "Ollama: Run 'ollama serve'"
        read -p "Enter local server URL (default: http://localhost:1234/v1): " LOCAL_URL
        OPENAI_API_KEY="not-needed"
        OPENAI_BASE_URL="${LOCAL_URL:-http://localhost:1234/v1}"
        read -p "Enter model name: " OPENAI_MODEL
        ;;
    4)
        read -p "Enter API Key: " OPENAI_API_KEY
        read -p "Enter Base URL: " OPENAI_BASE_URL
        read -p "Enter Model name: " OPENAI_MODEL
        ;;
esac

# Update .env file
sed -i.bak "s|^OPENAI_API_KEY=.*|OPENAI_API_KEY=$OPENAI_API_KEY|" .env
sed -i.bak "s|^OPENAI_BASE_URL=.*|OPENAI_BASE_URL=$OPENAI_BASE_URL|" .env
sed -i.bak "s|^OPENAI_MODEL=.*|OPENAI_MODEL=$OPENAI_MODEL|" .env

# Generate random session secret
SESSION_SECRET=$(openssl rand -base64 32 2>/dev/null || head -c 32 /dev/urandom | base64)
sed -i.bak "s|^SESSION_SECRET=.*|SESSION_SECRET=$SESSION_SECRET|" .env

# Clean up backup files
rm -f .env.bak

echo ""
echo -e "${GREEN}=========================================="
echo "  ✓ Setup Complete!"
echo "==========================================${NC}"
echo ""
echo "Your .env file has been configured with:"
echo "  - Google OAuth credentials"
echo "  - Microsoft OAuth credentials"
echo "  - AI provider settings"
echo "  - Secure session secret"
echo ""
echo "Next steps:"
echo "  1. Install dependencies: make install-deps"
echo "  2. Start the app: make dev"
echo "  3. Open http://localhost:3000"
echo ""
echo -e "${BLUE}Happy emailing! 🚀${NC}"
