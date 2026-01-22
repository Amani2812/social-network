# How to Restart Backend Server

## The Issue:
The backend server is currently running with the OLD code. The new like/dislike features won't work until you restart the server with the NEW code.

## Solution:

### Step 1: Stop the Current Backend Server
You need to find and stop the backend server that's running on port 8080.

**Option A - Using Task Manager (Windows):**
1. Press `Ctrl + Shift + Esc` to open Task Manager
2. Find the process named `server.exe` or `go.exe`
3. Right-click and select "End Task"

**Option B - Using Command Line:**
1. Open a new terminal
2. Run: `netstat -ano | findstr :8080`
3. Note the PID (Process ID) in the last column
4. Run: `taskkill /PID <PID> /F` (replace <PID> with the actual number)

**Option C - Using PowerShell:**
```powershell
Get-Process -Name "server" | Stop-Process -Force
```
OR
```powershell
Get-Process -Name "go" | Stop-Process -Force
```

### Step 2: Start the Backend Server with New Code
After stopping the old server:

1. Open a terminal in the backend directory
2. Run:
```bash
cd backend
go run server.go
```

3. You should see:
```
2026/01/22 XX:XX:XX Server starting on :8080
2026/01/22 XX:XX:XX WebSocket available at ws://localhost:8080/ws
```

### Step 3: Verify It's Working
1. Refresh your browser at `http://localhost:3001/dashboard`
2. Try clicking the 👍 or 👎 buttons on a post
3. The counts should update without errors!

## What Will Happen:

When the backend restarts:
1. ✅ New database tables will be created automatically (`post_reactions`, `comment_reactions`)
2. ✅ New API endpoints will be available (`/api/posts/react`, `/api/comments/react`)
3. ✅ Like/dislike buttons will work
4. ✅ Emoji picker already works (no backend needed)

## Quick Test:

After restarting, test these:
- ✅ Click 👍 on a post - count should increase
- ✅ Click 👍 again - count should decrease
- ✅ Click 👎 - dislike count should increase
- ✅ Open emoji picker in messages - should show 400+ emojis
- ✅ Click an emoji - should add to message

## If You Still See Errors:

1. Make sure the backend server actually restarted (check the terminal output)
2. Clear your browser cache (Ctrl + Shift + Delete)
3. Hard refresh the page (Ctrl + F5)
4. Check the browser console for any error messages

---

**Note:** The backend server must be restarted for the like/dislike feature to work. The emoji picker works immediately without any restart!
