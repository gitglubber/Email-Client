# AI-Powered Email & Calendar Client

An Outlook clone with AI-powered email sorting, auto-responses, spam filtering, and calendar management. Built with Go backend and React frontend.

## ⚡ Quick Start

**Get up and running in 2 minutes:**

```bash
./scripts/quickstart.sh
```

This automated script will:
- Guide you through OAuth setup for Google & Microsoft
- Configure your AI provider (OpenAI, Azure, or free local LLM)
- Install all dependencies
- Start the application

**Open http://localhost:3000 and sign in!**

---

## ✨ Features

- 🔐 **OAuth2 Authentication** - Sign in with Google and Microsoft accounts
- 📧 **Email Management** - Read, send, and organize emails via Gmail and Microsoft Graph APIs
- 🤖 **AI Email Sorting** - Automatically categorize emails using AI
- ✍️ **AI Auto-Responses** - Generate smart email replies
- 🛡️ **AI Spam Filtering** - Advanced spam detection using LLMs
- 📅 **Smart Calendar** - AI-powered calendar management and scheduling
- 🔌 **OpenAI Compatible** - Works with any OpenAI-compatible API endpoint (including **FREE local LLMs**)

## 🏗️ Architecture

- **Backend**: Go with Gin framework
- **Frontend**: React with TypeScript
- **Authentication**: OAuth2 for Google and Microsoft
- **APIs**: Gmail API, Microsoft Graph API
- **AI**: OpenAI-compatible endpoints (OpenAI, Azure OpenAI, local LLMs, etc.)
- **Database**: SQLite (easily swappable)

## 📖 Documentation

- **[QUICKSTART.md](./QUICKSTART.md)** - Get started in 5 minutes
- **[INSTALL_GUIDE.md](./INSTALL_GUIDE.md)** - Complete installation guide
- **[SETUP.md](./SETUP.md)** - Detailed setup instructions
- **[FEATURES.md](./FEATURES.md)** - Full feature list

## 🛠️ Setup Options

### Option 1: Automated Setup (Recommended)

Run the guided setup wizard:

```bash
./scripts/quickstart.sh
```

### Option 2: Automated OAuth Config Only

Just need help with OAuth credentials?

```bash
./scripts/setup-oauth.sh
```

### Option 3: Manual Setup

See [INSTALL_GUIDE.md](./INSTALL_GUIDE.md) for step-by-step instructions.

## 🚀 Running the Application

### Development Mode
```bash
make dev
```

Starts backend on http://localhost:8080 and frontend on http://localhost:3000

### Production Build
```bash
make build
./bin/email-client
```

## 🔧 Helpful Commands

```bash
# Check your environment and configuration
./scripts/check-env.sh

# Quick start with guided setup
./scripts/quickstart.sh

# Configure OAuth credentials only
./scripts/setup-oauth.sh

# Install all dependencies
make install-deps

# Run in development
make dev

# Build for production
make build

# Clean build artifacts
make clean
```

## 📋 Prerequisites

- Go 1.21 or higher
- Node.js 18+ and npm
- Google Account (for Gmail integration)
- Microsoft Account (optional, for Outlook integration)
- OpenAI API key OR local LLM (Ollama/LM Studio)

## 🤖 AI Provider Options

### OpenAI (Paid)
Fast, high-quality AI responses
- Get API key: https://platform.openai.com/api-keys
- Model: GPT-4 Turbo

### Local LLM (FREE!)
100% private, runs on your machine

**LM Studio** (Recommended for beginners):
1. Download from https://lmstudio.ai/
2. Download Mistral 7B model
3. Start local server
4. Configure in setup wizard

**Ollama** (For developers):
```bash
ollama pull mistral
ollama serve
```

### Azure OpenAI
Enterprise-ready with your own deployment

## 📸 Screenshots

![Login Page](docs/screenshots/login.png)
![Dashboard](docs/screenshots/dashboard.png)
![AI Reply](docs/screenshots/ai-reply.png)
![Calendar](docs/screenshots/calendar.png)

*Note: Screenshots to be added*

## API Endpoints

### Authentication
- `GET /auth/google` - Initiate Google OAuth flow
- `GET /auth/google/callback` - Google OAuth callback
- `GET /auth/microsoft` - Initiate Microsoft OAuth flow
- `GET /auth/microsoft/callback` - Microsoft OAuth callback
- `GET /auth/logout` - Logout current user
- `GET /auth/me` - Get current user info

### Email
- `GET /api/emails` - List emails
- `GET /api/emails/:id` - Get email details
- `POST /api/emails/send` - Send email
- `POST /api/emails/:id/reply` - Reply to email
- `POST /api/emails/:id/ai-reply` - Generate AI reply
- `POST /api/emails/:id/categorize` - AI categorize email
- `POST /api/emails/:id/spam-check` - Check if spam

### Calendar
- `GET /api/calendar/events` - List calendar events
- `POST /api/calendar/events` - Create event
- `POST /api/calendar/ai-schedule` - AI-powered scheduling

## Development

```bash
# Run backend with hot reload
go run cmd/server/main.go

# Run tests
go test ./...

# Build for production
go build -o bin/email-client cmd/server/main.go
```

## License

MIT
