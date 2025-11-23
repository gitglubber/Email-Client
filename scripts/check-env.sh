#!/bin/bash

# Environment Checker - Validates your setup before running

echo "=========================================="
echo "  Environment Check"
echo "=========================================="
echo ""

ERRORS=0
WARNINGS=0

# Check Go installation
echo -n "Checking Go installation... "
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | awk '{print $3}')
    echo "✓ Found $GO_VERSION"
else
    echo "✗ Go not found"
    echo "  Please install Go 1.21+ from https://go.dev/dl/"
    ERRORS=$((ERRORS+1))
fi

# Check Node.js installation
echo -n "Checking Node.js installation... "
if command -v node &> /dev/null; then
    NODE_VERSION=$(node --version)
    echo "✓ Found $NODE_VERSION"
else
    echo "✗ Node.js not found"
    echo "  Please install Node.js 18+ from https://nodejs.org/"
    ERRORS=$((ERRORS+1))
fi

# Check npm installation
echo -n "Checking npm installation... "
if command -v npm &> /dev/null; then
    NPM_VERSION=$(npm --version)
    echo "✓ Found v$NPM_VERSION"
else
    echo "✗ npm not found"
    ERRORS=$((ERRORS+1))
fi

# Check .env file
echo -n "Checking .env file... "
if [ -f .env ]; then
    echo "✓ Found"

    # Source .env
    source .env

    # Check Google OAuth
    echo -n "  Checking Google OAuth... "
    if [ -n "$GOOGLE_CLIENT_ID" ] && [ "$GOOGLE_CLIENT_ID" != "your-google-client-id" ]; then
        echo "✓ Configured"
    else
        echo "✗ Not configured"
        WARNINGS=$((WARNINGS+1))
    fi

    # Check Microsoft OAuth
    echo -n "  Checking Microsoft OAuth... "
    if [ -n "$MICROSOFT_CLIENT_ID" ] && [ "$MICROSOFT_CLIENT_ID" != "your-microsoft-client-id" ]; then
        echo "✓ Configured"
    else
        echo "✗ Not configured"
        WARNINGS=$((WARNINGS+1))
    fi

    # Check OpenAI
    echo -n "  Checking OpenAI API... "
    if [ -n "$OPENAI_API_KEY" ] && [ "$OPENAI_API_KEY" != "your-openai-api-key" ]; then
        echo "✓ Configured"
    else
        echo "⚠ Not configured (AI features won't work)"
        WARNINGS=$((WARNINGS+1))
    fi

else
    echo "✗ Not found"
    echo "  Run: cp .env.example .env"
    echo "  Then configure your credentials"
    ERRORS=$((ERRORS+1))
fi

# Check port availability
echo -n "Checking port 8080 (backend)... "
if lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null 2>&1; then
    echo "⚠ Port in use"
    WARNINGS=$((WARNINGS+1))
else
    echo "✓ Available"
fi

echo -n "Checking port 3000 (frontend)... "
if lsof -Pi :3000 -sTCP:LISTEN -t >/dev/null 2>&1; then
    echo "⚠ Port in use"
    WARNINGS=$((WARNINGS+1))
else
    echo "✓ Available"
fi

echo ""
echo "=========================================="
echo "  Results"
echo "=========================================="
echo ""

if [ $ERRORS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo "✓ All checks passed! Ready to run."
    echo ""
    echo "Start the application with:"
    echo "  ./scripts/quickstart.sh"
    echo "  or"
    echo "  make dev"
    exit 0
elif [ $ERRORS -gt 0 ]; then
    echo "✗ Found $ERRORS error(s) and $WARNINGS warning(s)"
    echo ""
    echo "Please fix the errors above before running."
    exit 1
else
    echo "⚠ Found $WARNINGS warning(s)"
    echo ""
    echo "You can still run the application, but some features may not work."
    echo ""
    echo "To fix warnings, run:"
    echo "  ./scripts/setup-oauth.sh"
    exit 0
fi
