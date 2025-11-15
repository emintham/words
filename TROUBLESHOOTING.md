# Troubleshooting Guide

## 403 / 404 Errors from Backend

If you're getting 403 or 404 errors when trying to use the application, follow these steps:

### 1. Check if the API Server is Running

The Go backend must be running for the application to work. Check if it's running:

```bash
# On Linux/Mac
lsof -i :9090

# On Windows
netstat -ano | findstr :9090
```

If nothing appears, the API server isn't running.

### 2. Start the API Server

From the project root directory:

```bash
# Build the API server (first time or after code changes)
go build -o api cmd/api/main.go

# Run the API server
./api
```

You should see output like:
```
2025/11/15 04:55:02 Starting server on port 9090
[GIN-debug] Listening and serving HTTP on :9090
```

**Keep this terminal window open** - the API server needs to stay running.

### 3. Start the Web Development Server

Open a **NEW terminal window** and run:

```bash
cd web
pnpm install  # Only needed first time or after package.json changes
pnpm run dev
```

You should see:
```
VITE v7.x.x  ready in xxx ms
➜  Local:   http://localhost:5173/
```

### 4. Access the Application

Open your browser to: **http://localhost:5173**

**Important:** Make sure you're using `localhost:5173` (the Vite dev server), NOT `localhost:9090` (the API server directly).

## Common Issues

### Issue: "Failed to communicate with server"

**Cause:** The API server isn't running or isn't accessible.

**Solution:**
1. Make sure the API server is running (see step 2 above)
2. Check that it's listening on port 9090
3. Verify nothing else is using port 9090

### Issue: CORS errors in browser console

**Cause:** Accessing the API directly instead of through the Vite proxy.

**Solution:**
- Always access via `http://localhost:5173`
- Don't access `http://localhost:9090` directly in your browser
- The Vite dev server proxies API requests correctly

### Issue: 404 on specific endpoints

**Cause:** The API routes might not be registered correctly.

**Solution:**
1. Restart the API server
2. Check the server logs for any errors
3. Verify the database file `words.db` exists

### Issue: Port already in use

**Error:** "address already in use" or "bind: address already in use"

**Solution:**
```bash
# Find what's using the port
lsof -i :9090  # or :5173

# Kill the process using the port
kill -9 <PID>

# Then restart the servers
```

## Quick Start Script

For convenience, you can use these commands in separate terminal windows:

**Terminal 1 (API Server):**
```bash
./start-api.sh
```

**Terminal 2 (Web Dev Server):**
```bash
cd web && pnpm run dev
```

## Still Having Issues?

1. Check both terminal windows for error messages
2. Open browser DevTools (F12) → Network tab to see exact error codes
3. Open browser DevTools (F12) → Console tab to see JavaScript errors
4. Verify your Go version: `go version` (should be 1.19+)
5. Verify your Node version: `node --version` (should be 18+)
