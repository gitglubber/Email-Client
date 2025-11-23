# Complete Installation Guide

This guide walks you through every step to get the AI Email Client running.

## Table of Contents

1. [System Requirements](#system-requirements)
2. [Installation Steps](#installation-steps)
3. [OAuth Configuration](#oauth-configuration)
4. [AI Provider Setup](#ai-provider-setup)
5. [Running the Application](#running-the-application)
6. [Verification](#verification)
7. [Troubleshooting](#troubleshooting)

## System Requirements

### Required Software

- **Go 1.21 or higher**
  - Download: https://go.dev/dl/
  - Verify: `go version`

- **Node.js 18 or higher**
  - Download: https://nodejs.org/
  - Verify: `node --version`

- **npm** (comes with Node.js)
  - Verify: `npm --version`

### Required Accounts

- **Google Account** (for Gmail integration)
- **Microsoft Account** (optional, for Outlook integration)
- **OpenAI Account** OR **Local LLM** (for AI features)

### System Recommendations

- 2GB RAM minimum
- 1GB free disk space
- Modern web browser (Chrome, Firefox, Safari, Edge)

## Installation Steps

### Step 1: Clone the Repository

If you haven't already:

```bash
git clone https://github.com/gitglubber/Email-Client.git
cd Email-Client
```

### Step 2: Run the Automated Setup

The easiest way to get started:

```bash
./scripts/quickstart.sh
```

This script will:
1. Check if you have all prerequisites
2. Guide you through OAuth setup
3. Configure AI provider
4. Install dependencies
5. Start the application

**That's it!** If the script runs successfully, skip to [Verification](#verification).

### Step 3: Manual Setup (Alternative)

If you prefer manual setup:

#### 3a. Create Environment File

```bash
cp .env.example .env
```

#### 3b. Configure Google OAuth

1. **Create Google Cloud Project**
   - Go to https://console.cloud.google.com/
   - Click "Select a project" → "New Project"
   - Name it "AI Email Client"
   - Click "Create"

2. **Enable Required APIs**
   - In the search bar, type "Gmail API" → Enable
   - Search "Google Calendar API" → Enable
   - Search "Google+ API" → Enable

3. **Configure OAuth Consent Screen**
   - Go to "APIs & Services" → "OAuth consent screen"
   - Choose "External" → "Create"
   - Fill in:
     - App name: "AI Email Client"
     - User support email: Your email
     - Developer contact: Your email
   - Click "Save and Continue"
   - Click "Add or Remove Scopes"
   - Add these scopes:
     - `.../auth/userinfo.email`
     - `.../auth/userinfo.profile`
     - `.../auth/gmail.readonly`
     - `.../auth/gmail.send`
     - `.../auth/gmail.modify`
     - `.../auth/calendar`
   - Click "Update" → "Save and Continue"
   - Add test users (your email) → "Save and Continue"

4. **Create OAuth Credentials**
   - Go to "Credentials" → "Create Credentials" → "OAuth client ID"
   - Choose "Web application"
   - Name: "AI Email Client"
   - Authorized redirect URIs: Add `http://localhost:8080/auth/google/callback`
   - Click "Create"
   - Copy the **Client ID** and **Client Secret**

5. **Update .env File**
   ```env
   GOOGLE_CLIENT_ID=your-client-id-here.apps.googleusercontent.com
   GOOGLE_CLIENT_SECRET=your-client-secret-here
   ```

#### 3c. Configure Microsoft OAuth (Optional)

1. **Register Application**
   - Go to https://portal.azure.com/
   - Navigate to "Azure Active Directory"
   - Click "App registrations" → "New registration"
   - Fill in:
     - Name: "AI Email Client"
     - Supported account types: "Personal Microsoft accounts only"
     - Redirect URI: Web - `http://localhost:8080/auth/microsoft/callback`
   - Click "Register"
   - Copy the **Application (client) ID**

2. **Create Client Secret**
   - Go to "Certificates & secrets"
   - Click "New client secret"
   - Description: "email-client"
   - Expires: Choose duration
   - Click "Add"
   - **Copy the Value immediately** (you won't see it again!)

3. **Add API Permissions**
   - Go to "API permissions"
   - Click "Add a permission" → "Microsoft Graph" → "Delegated permissions"
   - Search and add:
     - `User.Read`
     - `Mail.Read`
     - `Mail.Send`
     - `Mail.ReadWrite`
     - `Calendars.ReadWrite`
     - `offline_access`
   - Click "Add permissions"
   - Click "Grant admin consent" (if you have admin rights)

4. **Update .env File**
   ```env
   MICROSOFT_CLIENT_ID=your-application-id-here
   MICROSOFT_CLIENT_SECRET=your-client-secret-here
   ```

#### 3d. Configure AI Provider

Choose one option:

**Option A: OpenAI (Paid)**
```env
OPENAI_API_KEY=sk-your-key-here
OPENAI_BASE_URL=https://api.openai.com/v1
OPENAI_MODEL=gpt-4-turbo-preview
```

Get your key at: https://platform.openai.com/api-keys

**Option B: LM Studio (Free, Local)**
1. Download from https://lmstudio.ai/
2. Install and open LM Studio
3. Download a model:
   - Click "Search" tab
   - Search for "Mistral 7B Instruct" or "Llama 2"
   - Click download on a quantized version (Q4 recommended)
4. Start the server:
   - Click "Local Server" tab
   - Click "Start Server"
   - Note the URL (usually http://localhost:1234/v1)

```env
OPENAI_API_KEY=not-needed
OPENAI_BASE_URL=http://localhost:1234/v1
OPENAI_MODEL=local-model
```

**Option C: Ollama (Free, Local)**
1. Install Ollama from https://ollama.ai/
2. Pull a model:
   ```bash
   ollama pull mistral
   ```
3. Start the server:
   ```bash
   ollama serve
   ```

```env
OPENAI_API_KEY=not-needed
OPENAI_BASE_URL=http://localhost:11434/v1
OPENAI_MODEL=mistral
```

#### 3e. Generate Session Secret

```bash
# On Mac/Linux
openssl rand -base64 32

# Or use any random string generator
```

Update .env:
```env
SESSION_SECRET=your-random-secret-here
```

#### 3f. Install Dependencies

```bash
# Install Go dependencies
go mod download

# Install frontend dependencies
cd web
npm install
cd ..
```

## Running the Application

### Development Mode (Both servers)

```bash
make dev
```

This starts:
- Backend on http://localhost:8080
- Frontend on http://localhost:3000

### Run Separately

Terminal 1 (Backend):
```bash
make dev-backend
# or
go run cmd/server/main.go
```

Terminal 2 (Frontend):
```bash
make dev-frontend
# or
cd web && npm start
```

## Verification

### 1. Check Environment

```bash
./scripts/check-env.sh
```

This verifies all requirements are met.

### 2. Access the Application

Open your browser to: http://localhost:3000

You should see the login page with Google and Microsoft sign-in buttons.

### 3. Test Sign-In

1. Click "Sign in with Google"
2. Authorize the application
3. You should be redirected to the dashboard
4. Try viewing emails
5. Test AI features

### 4. Test AI Features

1. Click on any email
2. Click "AI Reply" - should generate a response
3. Click "AI Tools" → "Categorize" - should categorize the email
4. Click "AI Tools" → "Check Spam" - should check for spam

## Troubleshooting

### "Configuration Errors" on Startup

**Problem**: Server won't start, shows configuration errors.

**Solution**:
```bash
./scripts/setup-oauth.sh
```

### "Failed to exchange token"

**Problem**: OAuth redirect fails.

**Solutions**:
1. Verify redirect URI matches exactly:
   - Google: `http://localhost:8080/auth/google/callback`
   - Microsoft: `http://localhost:8080/auth/microsoft/callback`
2. Check you copied the entire Client ID and Secret
3. Ensure APIs are enabled (Google Cloud Console)
4. Clear browser cookies and try again

### AI Features Don't Work

**Problem**: AI reply/categorization fails.

**Solutions**:
1. Check your API key is valid
2. For local LLM, ensure server is running
3. Check backend logs for specific errors
4. Verify OPENAI_BASE_URL is correct

### Port Already in Use

**Problem**: Can't start server, port 8080 or 3000 in use.

**Solution**:
Change port in `.env`:
```env
PORT=8081
```

Update frontend proxy in `web/package.json`:
```json
"proxy": "http://localhost:8081"
```

### Network Errors

**Problem**: Can't fetch emails or calendar events.

**Solutions**:
1. Check internet connection
2. Verify OAuth credentials are correct
3. Check if Google/Microsoft services are accessible
4. Review backend logs for specific errors

### "Cannot find module" (Frontend)

**Problem**: Frontend build fails with missing modules.

**Solution**:
```bash
cd web
rm -rf node_modules package-lock.json
npm install
cd ..
```

### Go Build Fails

**Problem**: Cannot download Go modules.

**Solution**:
```bash
# Clear Go cache
go clean -modcache

# Re-download
go mod download

# Or use vendor
go mod vendor
go build -mod=vendor cmd/server/main.go
```

## Getting Help

If you're still stuck:

1. **Check Configuration**:
   ```bash
   ./scripts/check-env.sh
   ```

2. **Review Logs**:
   - Backend logs show in terminal
   - Frontend logs in browser console (F12)

3. **Read Documentation**:
   - [README.md](./README.md) - Overview
   - [SETUP.md](./SETUP.md) - Detailed setup
   - [FEATURES.md](./FEATURES.md) - Feature list
   - [QUICKSTART.md](./QUICKSTART.md) - Quick start guide

4. **Common Issues**:
   - Double-check all credentials
   - Ensure redirect URIs are exact
   - Verify all APIs are enabled
   - Check firewall isn't blocking localhost

## Next Steps

Once everything is working:

1. **Explore Features** - Try all AI capabilities
2. **Customize Prompts** - Edit `internal/ai/service.go`
3. **Add Themes** - Modify Tailwind config
4. **Deploy** - See production build instructions
5. **Contribute** - Submit PRs for improvements!

## Success Checklist

- [ ] Go and Node.js installed
- [ ] Google OAuth configured
- [ ] Microsoft OAuth configured (optional)
- [ ] AI provider configured
- [ ] Dependencies installed
- [ ] Server starts without errors
- [ ] Can sign in with Google
- [ ] Can view emails
- [ ] AI features work
- [ ] Calendar loads

Congratulations! You're all set up! 🎉
