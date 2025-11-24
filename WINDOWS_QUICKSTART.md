# Windows Quick Start Guide

Get the AI Email Client running on Windows in 5 minutes!

## Prerequisites

Install these first:

1. **Go 1.21+** - https://go.dev/dl/
   - Download the Windows installer
   - Run and follow the wizard
   - Verify: Open CMD and type `go version`

2. **Node.js 18+** - https://nodejs.org/
   - Download the Windows installer (LTS version)
   - Run and follow the wizard
   - Verify: Open CMD and type `node --version`

3. **Git for Windows** (if not already installed)
   - https://git-scm.com/download/win

## Quick Start

### Option 1: One-Command Setup (Easiest)

Open **Command Prompt** or **PowerShell** in the project directory and run:

```cmd
scripts\quickstart.bat
```

Or in PowerShell:
```powershell
.\scripts\quickstart.ps1
```

This will:
- Guide you through OAuth setup
- Install all dependencies
- Start the application

### Option 2: Step-by-Step

#### Step 1: Check Your Environment

```cmd
scripts\check-env.bat
```

This verifies:
- Go and Node.js are installed
- OAuth credentials are configured
- Ports are available

#### Step 2: Configure OAuth

```cmd
scripts\setup-oauth.bat
```

Follow the interactive wizard to set up:
- Google OAuth (required)
- Microsoft OAuth (optional)
- AI provider (OpenAI or free local LLM)

#### Step 3: Install Dependencies

```cmd
go mod download
cd web
npm install
cd ..
```

#### Step 4: Run the Application

**Option A: Using Make (if you have it)**
```cmd
make dev
```

**Option B: Manual Start**

Terminal 1 - Backend:
```cmd
go run cmd\server\main.go
```

Terminal 2 - Frontend:
```cmd
cd web
npm start
```

#### Step 5: Open the Application

Open your browser to: **http://localhost:3000**

## Using PowerShell

If you prefer PowerShell, all scripts have `.ps1` versions:

```powershell
# Setup OAuth
.\scripts\setup-oauth.ps1

# Check environment
.\scripts\check-env.ps1

# Quick start
.\scripts\quickstart.ps1
```

**Note:** You may need to enable script execution:
```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

## Troubleshooting

### "go: command not found"
- Go is not installed or not in PATH
- Solution: Reinstall Go from https://go.dev/dl/
- Make sure to restart your terminal after installation

### "node: command not found"
- Node.js is not installed or not in PATH
- Solution: Reinstall Node.js from https://nodejs.org/
- Restart your terminal after installation

### "Cannot run scripts" (PowerShell)
Error: `execution of scripts is disabled on this system`

Solution:
```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### Port Already in Use

If port 8080 or 3000 is in use:

1. Find and kill the process:
```cmd
netstat -ano | findstr :8080
taskkill /PID <PID> /F
```

2. Or change the port in `.env`:
```
PORT=8081
```

### Network Errors

If `go mod download` fails:

```cmd
# Set Go proxy
set GOPROXY=https://proxy.golang.org,direct

# Then retry
go mod download
```

### npm Install Fails

```cmd
cd web
rmdir /s /q node_modules
del package-lock.json
npm install
cd ..
```

## AI Provider Setup

### Free Option: LM Studio (Recommended for Windows)

1. Download LM Studio from https://lmstudio.ai/
2. Install and open LM Studio
3. Click "Search" tab
4. Search for "Mistral 7B Instruct"
5. Download a Q4 quantized version
6. Go to "Local Server" tab
7. Click "Start Server"
8. Note the URL (usually `http://localhost:1234/v1`)

In the setup wizard, choose option 3 (Local LLM) and enter the URL.

### Paid Option: OpenAI

1. Go to https://platform.openai.com/api-keys
2. Create an API key
3. In the setup wizard, choose option 1 and enter your key

## Development Tips

### Hot Reload Development

For the best development experience, run in two separate terminals:

Terminal 1 (Backend with auto-restart):
```cmd
go install github.com/cosmtrek/air@latest
air
```

Terminal 2 (Frontend):
```cmd
cd web
npm start
```

### Building for Production

```cmd
# Build backend
go build -o bin\email-client.exe cmd\server\main.go

# Build frontend
cd web
npm run build
cd ..

# Run production build
bin\email-client.exe
```

## What's Next?

Once the application is running:

1. **Sign in** - Choose Google or Microsoft
2. **View emails** - Your inbox loads automatically
3. **Try AI features**:
   - Click "AI Reply" on any email
   - Use "AI Tools" → "Categorize"
   - Use "AI Tools" → "Check Spam"
4. **Switch to Calendar** - View and manage events
5. **AI Scheduling** - Click "AI Schedule" for smart suggestions

## Getting Help

- **Check environment**: `scripts\check-env.bat`
- **View logs**: Backend logs appear in the terminal
- **Browser console**: Press F12 to see frontend logs
- **Documentation**: Read `QUICKSTART.md` and `INSTALL_GUIDE.md`

## Common Commands

```cmd
REM Check everything is installed
scripts\check-env.bat

REM Setup OAuth credentials
scripts\setup-oauth.bat

REM Start application
scripts\quickstart.bat

REM Or manually
go run cmd\server\main.go
cd web && npm start
```

## Success! 🎉

If you see:
- Backend running on http://localhost:8080
- Frontend on http://localhost:3000
- Login page with Google/Microsoft buttons

**You're all set!** Sign in and start using your AI-powered email client!
