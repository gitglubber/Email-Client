# Quick Start Guide

Get up and running in 5 minutes!

## Prerequisites

- **Go 1.21+** - [Download here](https://go.dev/dl/)
- **Node.js 18+** - [Download here](https://nodejs.org/)
- **Google Account** - For Gmail/Calendar access
- **Microsoft Account** - For Outlook/Calendar access (optional)
- **OpenAI API Key** - Or use local LLM like Ollama/LM Studio

## Option 1: Automated Setup (Recommended)

Run the setup wizard that guides you through everything:

```bash
./scripts/quickstart.sh
```

This will:
1. Guide you through OAuth setup for Google & Microsoft
2. Configure your AI provider (OpenAI, Azure, or local)
3. Install all dependencies
4. Start the application

## Option 2: Manual Setup

### Step 1: OAuth Credentials

Run the OAuth setup wizard:

```bash
./scripts/setup-oauth.sh
```

This interactive script will:
- Guide you through Google Cloud Console setup
- Guide you through Microsoft Azure Portal setup
- Configure your OpenAI/AI settings
- Generate secure session secrets

### Step 2: Install Dependencies

```bash
# Install Go dependencies
go mod download

# Install frontend dependencies
cd web && npm install && cd ..
```

### Step 3: Run the Application

```bash
# Start both backend and frontend
make dev
```

Or run separately:

```bash
# Terminal 1 - Backend
make dev-backend

# Terminal 2 - Frontend
make dev-frontend
```

### Step 4: Access the Application

Open your browser to: **http://localhost:3000**

## Verify Your Setup

Check if everything is configured correctly:

```bash
./scripts/check-env.sh
```

This will verify:
- Go and Node.js are installed
- .env file exists and is configured
- OAuth credentials are set
- Ports are available

## Google OAuth Setup (Detailed)

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project
3. Enable these APIs:
   - Gmail API
   - Google Calendar API
   - Google+ API
4. Go to **Credentials** → **Create Credentials** → **OAuth 2.0 Client ID**
5. Configure OAuth consent screen:
   - User Type: **External**
   - Add scopes: email, profile, gmail.readonly, gmail.send, gmail.modify, calendar
6. Create OAuth Client ID:
   - Application type: **Web application**
   - Authorized redirect URIs: `http://localhost:8080/auth/google/callback`
7. Copy **Client ID** and **Client Secret**

## Microsoft OAuth Setup (Detailed)

1. Go to [Azure Portal](https://portal.azure.com/)
2. Navigate to **Azure Active Directory** → **App registrations** → **New registration**
3. Configure:
   - Name: `AI Email Client`
   - Supported account types: **Personal Microsoft accounts**
   - Redirect URI: `http://localhost:8080/auth/microsoft/callback`
4. After creation, note the **Application (client) ID**
5. Go to **Certificates & secrets** → **New client secret**
   - Create and copy the secret value
6. Go to **API permissions** → **Add permission** → **Microsoft Graph**
   - Add these delegated permissions:
     - User.Read
     - Mail.Read
     - Mail.Send
     - Mail.ReadWrite
     - Calendars.ReadWrite
     - offline_access
7. Click **Grant admin consent**

## OpenAI Setup Options

### Option A: OpenAI

```env
OPENAI_API_KEY=sk-your-key-here
OPENAI_BASE_URL=https://api.openai.com/v1
OPENAI_MODEL=gpt-4-turbo-preview
```

Get your API key: https://platform.openai.com/api-keys

### Option B: Local LLM (Free!)

Using **LM Studio**:
1. Download from https://lmstudio.ai/
2. Download a model (e.g., Mistral 7B Instruct)
3. Go to **Settings** → **Server** → **Start Server**
4. Use these settings:

```env
OPENAI_API_KEY=not-needed
OPENAI_BASE_URL=http://localhost:1234/v1
OPENAI_MODEL=local-model
```

Using **Ollama**:
1. Install from https://ollama.ai/
2. Run: `ollama pull mistral`
3. Run: `ollama serve`
4. Use these settings:

```env
OPENAI_API_KEY=not-needed
OPENAI_BASE_URL=http://localhost:11434/v1
OPENAI_MODEL=mistral
```

### Option C: Azure OpenAI

```env
OPENAI_API_KEY=your-azure-key
OPENAI_BASE_URL=https://YOUR-RESOURCE.openai.azure.com/
OPENAI_MODEL=gpt-4
```

## Troubleshooting

### "OAuth credentials not configured"

Run the setup wizard:
```bash
./scripts/setup-oauth.sh
```

### "Port 8080 already in use"

Change the port in `.env`:
```env
PORT=8081
```

And update the frontend proxy in `web/package.json`:
```json
"proxy": "http://localhost:8081"
```

### "Failed to exchange token"

1. Verify your redirect URIs match exactly:
   - Google: `http://localhost:8080/auth/google/callback`
   - Microsoft: `http://localhost:8080/auth/microsoft/callback`
2. Check that you've enabled all required APIs/permissions
3. Clear browser cookies and try again

### AI features not working

1. Verify your API key is valid
2. Check the base URL is correct
3. For local LLMs, ensure the server is running
4. Check backend logs: `go run cmd/server/main.go`

### "Connection refused" errors

Make sure:
1. Backend is running on port 8080
2. Frontend proxy is configured correctly
3. No firewall blocking localhost connections

## What's Next?

Once running:

1. **Sign in** with Google or Microsoft
2. **View emails** - Your inbox loads automatically
3. **Try AI features**:
   - Click "AI Reply" on any email
   - Use "AI Tools" → "Categorize" to sort emails
   - Use "AI Tools" → "Check Spam" for spam detection
4. **Manage calendar** - Switch to Calendar view
5. **Use AI Scheduling** - Click "AI Schedule" for smart suggestions

## Production Deployment

For production:

1. Build the application:
   ```bash
   make build
   ```

2. Update OAuth redirect URIs to your domain:
   ```
   https://yourdomain.com/auth/google/callback
   https://yourdomain.com/auth/microsoft/callback
   ```

3. Set production environment variables:
   ```env
   FRONTEND_URL=https://yourdomain.com
   SESSION_SECRET=<strong-random-value>
   ```

4. Serve the frontend build from `web/build/`

5. Run the backend:
   ```bash
   ./bin/email-client
   ```

## Getting Help

- Run environment check: `./scripts/check-env.sh`
- Check detailed setup: Read `SETUP.md`
- View all features: Read `FEATURES.md`
- Review main docs: Read `README.md`

## Common Commands

```bash
# Quick start with wizard
./scripts/quickstart.sh

# Setup OAuth credentials
./scripts/setup-oauth.sh

# Check environment
./scripts/check-env.sh

# Install dependencies
make install-deps

# Run in development
make dev

# Build for production
make build

# Clean build artifacts
make clean
```

Happy emailing! 🚀
