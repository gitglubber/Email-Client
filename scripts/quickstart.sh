#!/bin/bash

# Quick Start Script - Sets up and runs the application

set -e

echo "=========================================="
echo "  AI Email Client - Quick Start"
echo "=========================================="
echo ""

# Check if .env exists and has credentials
if [ ! -f .env ]; then
    echo "No .env file found. Running OAuth setup..."
    echo ""
    ./scripts/setup-oauth.sh
else
    # Check if credentials are set
    source .env
    if [ "$GOOGLE_CLIENT_ID" = "your-google-client-id" ] || [ "$MICROSOFT_CLIENT_ID" = "your-microsoft-client-id" ]; then
        echo "OAuth credentials not configured."
        echo ""
        read -p "Would you like to run the setup wizard? (y/n) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            ./scripts/setup-oauth.sh
        else
            echo "Please configure .env manually or run: ./scripts/setup-oauth.sh"
            exit 1
        fi
    fi
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.21+ from https://go.dev/dl/"
    exit 1
fi

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    echo "Error: Node.js is not installed. Please install Node.js 18+ from https://nodejs.org/"
    exit 1
fi

echo ""
echo "Installing dependencies..."
echo ""

# Install Go dependencies
echo "Installing Go dependencies..."
go mod download

# Install frontend dependencies
echo "Installing frontend dependencies..."
cd web
npm install --silent
cd ..

echo ""
echo "=========================================="
echo "  Starting Application"
echo "=========================================="
echo ""
echo "Backend: http://localhost:8080"
echo "Frontend: http://localhost:3000"
echo ""
echo "Press Ctrl+C to stop both servers"
echo ""

# Start both servers
make dev
