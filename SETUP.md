# Setup Guide

Complete setup instructions for the AI-Powered Email & Calendar Client.

## Prerequisites

- Go 1.21 or higher
- Node.js 18+ and npm
- Google Cloud Console account
- Microsoft Azure account
- OpenAI API key (or compatible endpoint like Azure OpenAI, LM Studio, Ollama)

## Step 1: Clone and Setup Environment

```bash
cd Email-Client
cp .env.example .env
```

Edit `.env` file with your credentials (see configuration steps below).

## Step 2: Google OAuth Configuration

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing one
3. Enable the following APIs:
   - Gmail API
   - Google Calendar API
   - Google+ API (for user info)

4. Go to "Credentials" → "Create Credentials" → "OAuth 2.0 Client ID"
5. Configure OAuth consent screen:
   - User Type: External
   - Add scopes: email, profile, Gmail, Calendar
6. Create OAuth 2.0 Client ID:
   - Application type: Web application
   - Authorized redirect URIs: `http://localhost:8080/auth/google/callback`
   - For production, add your domain: `https://yourdomain.com/auth/google/callback`

7. Copy Client ID and Client Secret to `.env`:
```env
GOOGLE_CLIENT_ID=your-google-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback
```

## Step 3: Microsoft OAuth Configuration

1. Go to [Azure Portal](https://portal.azure.com/)
2. Navigate to "Azure Active Directory" → "App registrations" → "New registration"
3. Register application:
   - Name: AI Email Client
   - Supported account types: Accounts in any organizational directory and personal Microsoft accounts
   - Redirect URI: Web - `http://localhost:8080/auth/microsoft/callback`

4. After registration, note the "Application (client) ID"

5. Go to "Certificates & secrets" → "New client secret"
   - Description: email-client-secret
   - Expires: 24 months
   - Copy the secret value (you'll only see it once!)

6. Go to "API permissions" → "Add a permission" → "Microsoft Graph" → "Delegated permissions"
   - Add these permissions:
     - User.Read
     - Mail.Read
     - Mail.Send
     - Mail.ReadWrite
     - Calendars.ReadWrite
     - offline_access

7. Click "Grant admin consent" (if you have admin rights)

8. Update `.env`:
```env
MICROSOFT_CLIENT_ID=your-microsoft-client-id
MICROSOFT_CLIENT_SECRET=your-microsoft-client-secret
MICROSOFT_REDIRECT_URL=http://localhost:8080/auth/microsoft/callback
```

## Step 4: OpenAI Configuration

### Option A: OpenAI
```env
OPENAI_API_KEY=sk-your-openai-api-key
OPENAI_BASE_URL=https://api.openai.com/v1
OPENAI_MODEL=gpt-4-turbo-preview
```

### Option B: Azure OpenAI
```env
OPENAI_API_KEY=your-azure-openai-key
OPENAI_BASE_URL=https://YOUR_RESOURCE.openai.azure.com/
OPENAI_MODEL=gpt-4
```

### Option C: Local LLM (LM Studio, Ollama, etc.)
```env
OPENAI_API_KEY=not-needed
OPENAI_BASE_URL=http://localhost:1234/v1
OPENAI_MODEL=local-model
```

For LM Studio:
1. Download and install [LM Studio](https://lmstudio.ai/)
2. Download a model (recommended: Mistral 7B Instruct or similar)
3. Start the local server in LM Studio (Settings → Server → Start Server)
4. Use the URL shown (usually http://localhost:1234/v1)

## Step 5: Install Dependencies

```bash
# Install Go dependencies
go mod download

# Install frontend dependencies
cd web
npm install
cd ..
```

## Step 6: Run the Application

### Development Mode (Both Backend and Frontend)
```bash
make dev
```

This starts:
- Backend server on http://localhost:8080
- Frontend on http://localhost:3000

### Run Backend Only
```bash
make dev-backend
# or
go run cmd/server/main.go
```

### Run Frontend Only
```bash
make dev-frontend
# or
cd web && npm start
```

## Step 7: Access the Application

1. Open your browser to http://localhost:3000
2. Click "Sign in with Google" or "Sign in with Microsoft"
3. Authorize the application
4. Start using your AI-powered email client!

## Production Build

```bash
# Build both backend and frontend
make build

# Run the production server
./bin/email-client
```

The frontend build will be in `web/build/`. You can serve it with any static file server or configure the Go backend to serve it.

## Features to Try

### Email Features
1. **View Inbox** - See all your emails with modern UI
2. **AI Categorization** - Click "AI Tools" → "Categorize" on any email
3. **AI Reply** - Click "AI Reply" to generate smart responses
4. **Spam Detection** - Click "AI Tools" → "Check Spam"
5. **Manual Reply** - Click "Reply" to compose manually

### Calendar Features
1. **View Events** - See upcoming calendar events
2. **Create Event** - Click "New Event" to create manually
3. **AI Scheduling** - Click "AI Schedule" for smart scheduling suggestions

## Troubleshooting

### "Failed to load data" or Authentication Issues
- Check that your OAuth credentials are correct
- Verify redirect URLs match exactly in Google/Microsoft consoles
- Make sure you've enabled all required APIs
- Clear browser cookies and try again

### AI Features Not Working
- Verify your OpenAI API key is valid
- Check the base URL is correct
- For local LLMs, ensure the server is running
- Check backend logs for specific error messages

### CORS Errors
- Ensure `FRONTEND_URL` in `.env` matches your frontend URL
- For production, update CORS settings in `internal/api/router.go`

### Port Already in Use
- Change `PORT` in `.env` to a different port
- Update frontend proxy in `web/package.json` to match

## Environment Variables Reference

| Variable | Description | Default |
|----------|-------------|---------|
| PORT | Backend server port | 8080 |
| FRONTEND_URL | Frontend URL for CORS | http://localhost:3000 |
| GOOGLE_CLIENT_ID | Google OAuth client ID | Required |
| GOOGLE_CLIENT_SECRET | Google OAuth secret | Required |
| GOOGLE_REDIRECT_URL | Google OAuth redirect | http://localhost:8080/auth/google/callback |
| MICROSOFT_CLIENT_ID | Microsoft OAuth client ID | Required |
| MICROSOFT_CLIENT_SECRET | Microsoft OAuth secret | Required |
| MICROSOFT_REDIRECT_URL | Microsoft OAuth redirect | http://localhost:8080/auth/microsoft/callback |
| OPENAI_API_KEY | OpenAI API key | Required |
| OPENAI_BASE_URL | API endpoint URL | https://api.openai.com/v1 |
| OPENAI_MODEL | Model to use | gpt-4-turbo-preview |
| SESSION_SECRET | Secret for sessions | change-me-in-production |

## Next Steps

- Customize AI prompts in `internal/ai/service.go`
- Add more email filters and labels
- Implement email search
- Add support for attachments
- Set up database for persistent storage
- Deploy to production server

## Support

For issues and questions:
- Check the main [README.md](./README.md)
- Review Go backend logs
- Check browser console for frontend errors
- Verify all environment variables are set correctly
