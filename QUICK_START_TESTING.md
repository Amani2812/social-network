# Quick Start - Real-Time Notifications Testing

## 🚀 Start the Application

### Terminal 1 - Backend
```bash
cd backend
go run server.go
```
**Expected output:**
```
Server starting on :8080
Client registered: 1
Client registered: 2
```

### Terminal 2 - Frontend
```bash
cd frontend
npm run dev
```
**Expected output:**
```
- ready started server on 0.0.0.0:3000
- Local: http://localhost:3000
```

---

## 🧪 Quick Test (2 Minutes)

### Setup
1. **Chrome Browser** → `http://localhost:3000`
   - Register/Login as: `alice@test.com` / `password123`
   
2. **Firefox Browser** → `http://localhost:3000`
   - Register/Login as: `bob@test.com` / `password123`

### Test 1: Group Join Request Notification

**In Chrome (Alice):**
1. Click **"Groups"** in navigation
2. Click **"Create Group"**
3. Enter:
   - Title: `Test Group`
   - Description: `Testing notifications`
4. Click **"Create"**
5. **👀 WATCH THE TOP-RIGHT CORNER** for notifications

**In Firefox (Bob):**
1. Click **"Groups"**
2. Find "Test Group"
3. Click **"Request to Join"**

**✅ Expected Result in Chrome:**
- 🎉 Toast notification appears: "Bob Member wants to join Test Group"
- 🔔 Bell badge shows "1"
- No page refresh needed!

---

### Test 2: Event Creation Notification

**In Chrome (Alice):**
1. Click notification bell 🔔
2. Find Bob's join request
3. Click **"Accept"** ✓

**In Firefox (Bob):**
1. Refresh Groups page
2. Click on "Test Group"
3. **👀 WATCH THE TOP-RIGHT CORNER**

**In Chrome (Alice):**
1. In "Test Group", click **"Events"** tab
2. Click **"+ Create Event"**
3. Enter:
   - Title: `Team Meeting`
   - Date: Any future date
4. Click **"Create Event"**

**✅ Expected Result in Firefox:**
- 🎉 Toast notification appears: "New event in Test Group: Team Meeting"
- 🔔 Bell badge updates
- No page refresh needed!

---

## 🔍 Verify It's Working

### Check Browser Console (F12)
**You should see:**
```
✅ Dashboard WebSocket connected
📥 Notification received: {type: "notification", content: "..."}
```

### Check Backend Terminal
**You should see:**
```
Client registered: 1
Client registered: 2
Sent ping to client 1
Sent ping to client 2
```

---

## ❌ Troubleshooting

### Problem: No notification appears

**Solution 1:** Check WebSocket connection
- Open Browser Console (F12)
- Look for: `✅ WebSocket connected`
- If not connected, restart backend

**Solution 2:** Check backend is running
```bash
# Should show backend process
netstat -ano | findstr :8080
```

**Solution 3:** Clear browser cache
- Press `Ctrl + Shift + Delete`
- Clear cache and reload

### Problem: "Connection refused"

**Check:**
1. Backend running on port 8080?
2. Frontend running on port 3000?
3. Firewall blocking connections?

**Fix:**
```bash
# Restart backend
cd backend
go run server.go

# Restart frontend
cd frontend
npm run dev
```

---

## 📊 Success Indicators

### ✅ Everything Working If:
- [ ] Toast notifications appear instantly
- [ ] Bell badge updates without refresh
- [ ] Console shows WebSocket connected
- [ ] Backend logs show client registrations
- [ ] Notifications persist in database

### ❌ Something Wrong If:
- [ ] Need to refresh to see notifications
- [ ] Console shows WebSocket errors
- [ ] No toast notifications appear
- [ ] Bell badge doesn't update

---

## 🎯 What to Test

### Core Features:
1. **Group Join Requests** ✅
   - Create group → Request to join → Notification appears

2. **Event Creation** ✅
   - Create event → All members notified

3. **Multiple Notifications** ✅
   - Send multiple requests → All appear

4. **Notification Persistence** ✅
   - Close browser → Reopen → Notifications still there

5. **Real-Time Updates** ✅
   - No refresh needed → Instant updates

---

## 📝 Test Checklist

```
□ Backend starts without errors
□ Frontend starts without errors
□ Can register/login two users
□ Can create a group
□ Can request to join group
□ Admin receives real-time notification
□ Can accept join request
□ Can create event
□ Members receive real-time notification
□ Toast notifications appear
□ Bell badge updates
□ Notifications page shows all notifications
□ WebSocket stays connected
□ Auto-reconnects if disconnected
```

---

## 🎉 Demo Script

**For showing to others:**

```
1. "Let me show you real-time notifications..."

2. [Open Chrome] "This is Alice, the group admin"
   - Create group "Demo Group"

3. [Open Firefox] "This is Bob, who wants to join"
   - Request to join "Demo Group"

4. [Point to Chrome] "Watch Alice's screen..."
   - 🎉 Notification appears instantly!

5. [Chrome] Accept Bob's request

6. [Chrome] Create event "Team Meeting"

7. [Point to Firefox] "Watch Bob's screen..."
   - 🎉 Notification appears instantly!

8. "No page refresh needed - it's all real-time!"
```

---

## 📚 More Information

- **Detailed Testing**: See `MULTI_BROWSER_TESTING_GUIDE.md`
- **Implementation Details**: See `REALTIME_NOTIFICATIONS_COMPLETE.md`
- **Troubleshooting**: See `MULTI_BROWSER_TESTING_GUIDE.md` → Troubleshooting section

---

## 🚨 Important Notes

1. **Keep Both Browsers Open**: Notifications only work when browser is open
2. **WebSocket Required**: Modern browser needed (Chrome, Firefox, Edge, Safari)
3. **Same Network**: Both browsers should access same localhost
4. **Incognito Mode**: Use for testing with same browser (Chrome Incognito + Chrome Normal)

---

## ✨ Features You Can Test

- ✅ Real-time group join request notifications
- ✅ Real-time event creation notifications
- ✅ Toast notifications (auto-dismiss after 5s)
- ✅ Notification bell badge
- ✅ Notification persistence
- ✅ Click to navigate
- ✅ Mark as read/unread
- ✅ Auto-reconnection

---

**Ready to test? Let's go! 🚀**
