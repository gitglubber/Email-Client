# AI-Powered Email & Calendar Client

An Outlook clone with AI-powered email sorting, auto-responses, spam filtering, and calendar management. Built with Go backend and React frontend.

## Features

- 🔐 **OAuth2 Authentication** - Sign in with Google and Microsoft accounts
- 📧 **Email Management** - Read, send, and organize emails via Gmail and Microsoft Graph APIs
- 🤖 **AI Email Sorting** - Automatically categorize emails using AI
- ✍️ **AI Auto-Responses** - Generate smart email replies
- 🛡️ **AI Spam Filtering** - Advanced spam detection using LLMs
- 📅 **Smart Calendar** - AI-powered calendar management and scheduling
- 🔌 **OpenAI Compatible** - Works with any OpenAI-compatible API endpoint

## Architecture

- **Backend**: Go with Gin framework
- **Frontend**: React with TypeScript
- **Authentication**: OAuth2 for Google and Microsoft
- **APIs**: Gmail API, Microsoft Graph API
- **AI**: OpenAI-compatible endpoints (OpenAI, Azure OpenAI, local LLMs, etc.)
- **Database**: SQLite (easily swappable)

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Node.js 18+ and npm
- Google Cloud Console project with Gmail API enabled
- Microsoft Azure AD app registration
- OpenAI API key (or compatible endpoint)

### Setup

1. Clone the repository and copy environment file:
```bash
cp .env.example .env
```

2. Configure your `.env` file with your credentials

3. Run the backend:
```bash
go mod download
go run cmd/server/main.go
```

4. Run the frontend:
```bash
cd web
npm install
npm start
```

5. Open http://localhost:3000

## Configuration

### Google OAuth Setup

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project
3. Enable Gmail API
4. Create OAuth 2.0 credentials
5. Add authorized redirect URI: `http://localhost:8080/auth/google/callback`
6. Copy Client ID and Secret to `.env`

### Microsoft OAuth Setup

1. Go to [Azure Portal](https://portal.azure.com/)
2. Register a new application
3. Add permissions: Mail.Read, Mail.Send, Calendars.ReadWrite
4. Add redirect URI: `http://localhost:8080/auth/microsoft/callback`
5. Copy Application (client) ID and secret to `.env`

### OpenAI Compatible Endpoints

The app works with any OpenAI-compatible API:
- **OpenAI**: `https://api.openai.com/v1`
- **Azure OpenAI**: `https://YOUR_RESOURCE.openai.azure.com/`
- **Local (LM Studio, Ollama)**: `http://localhost:1234/v1`
- **Other providers**: Any API following OpenAI's spec

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
