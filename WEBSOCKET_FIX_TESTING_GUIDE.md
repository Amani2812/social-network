# WebSocket Fix - Testing Guide

## Overview
This guide provides step-by-step instructions for testing the WebSocket error fixes implemented across all 4 pages of the social network application.

## Prerequisites
- Backend server running on `http://localhost:8080`
- Frontend running on `http://localhost:3000`
- Browser with Developer Tools (Chrome recommended)
- At least 2 user accounts for testing

## Testing Steps

### 1. Dashboard Page Testing

#### Basic Connection Test
1. Open browser and navigate to `http://localhost:3000/dashboard`
2. Open Developer Tools (F12) → Console tab
3. Look for: `✅ Dashboard WebSocket connected`
4. Verify no error messages appear

#### Real-time Notifications Test
1. Keep dashboard open in one browser/tab
2. Open another browser/tab with a different user
3. Have the second user follow you or send a notification
4. Verify toast notification appears on dashboard
5. Check console for: `📥 Received:` messages

#### Reconnection Test
1. With dashboard open, stop the backend server
2. Check console for: `🔌 WebSocket closed:`
3. Check console for: `⏳ Reconnecting in X seconds...`
4. Restart backend server
5. Verify: `✅ Dashboard WebSocket connected` appears
6. Verify reconnection counter resets

#### Cleanup Test
1. Navigate to dashboard
2. Wait for WebSocket to connect
3. Navigate to another page (e.g., `/profile`)
4. Check console - should see proper cleanup
5. No error messages should appear

---

### 2. Notifications Page Testing

#### Basic Connection Test
1. Navigate to `http://localhost:3000/notifications`
2. Open Console
3. Look for: `✅ Notifications WebSocket connected`
4. Verify no errors

#### Real-time Updates Test
1. Keep notifications page open
2. Have another user send you a follow request
3. Verify notification appears in real-time
4. Check console for WebSocket messages

#### Reconnection Test
1. With notifications page open, stop backend
2. Observe reconnection attempts in console
3. Restart backend
4. Verify successful reconnection

---

### 3. Messages Page Testing

#### Basic Connection Test
1. Navigate to `http://localhost:3000/messages`
2. Open Console
3. Look for: `✅ WebSocket Connected`
4. Check connection status indicator (green dot = connected)

#### Real-time Messaging Test
1. Open messages page with User A
2. Open messages page with User B in another browser
3. Send message from User A to User B
4. Verify message appears instantly for User B
5. Send reply from User B
6. Verify reply appears instantly for User A

#### Connection Status Indicator Test
1. Observe connection status at top of page
2. Should show: "Connected" with green dot
3. Stop backend server
4. Should show: "Disconnected" with red dot
5. Restart backend
6. Should show: "Connecting..." then "Connected"

#### Message Queue Test
1. Stop backend server
2. Try sending a message
3. Message should be queued
4. Restart backend
5. Queued message should be sent automatically

---

### 4. Profile Page Testing

#### Basic Connection Test
1. Navigate to any user profile (e.g., `http://localhost:3000/profile/1`)
2. Open Console
3. Look for: `✅ Profile WebSocket connected`
4. Verify no errors

#### Follow Status Updates Test
1. Open User A's profile
2. In another browser, have User B accept User A's follow request
3. Verify follow status updates in real-time on User A's profile
4. Check console for: `📥 Profile received:` messages

#### Reconnection Test
1. With profile page open, stop backend
2. Observe reconnection attempts
3. Restart backend
4. Verify successful reconnection

---

## Advanced Testing

### Exponential Backoff Test
1. Stop backend server
2. Open any page with WebSocket
3. Observe console logs for reconnection timing:
   - Attempt 1: ~1 second
   - Attempt 2: ~2 seconds
   - Attempt 3: ~4 seconds
   - Attempt 4: ~8 seconds
   - Attempt 5: ~16 seconds
   - Attempt 6+: ~30 seconds (capped)

### Max Reconnection Attempts Test
1. Stop backend server
2. Open any page with WebSocket
3. Wait for 10 reconnection attempts
4. After 10 attempts, should see: `❌ Max reconnection attempts reached. Please refresh the page.`
5. No more reconnection attempts should occur

### Multiple Tabs Test
1. Open dashboard in 3 different tabs
2. All should connect successfully
3. Stop backend
4. All should attempt to reconnect
5. Restart backend
6. All should reconnect successfully

### Memory Leak Test
1. Open dashboard
2. Wait for WebSocket to connect
3. Navigate to notifications
4. Navigate to messages
5. Navigate to profile
6. Navigate back to dashboard
7. Repeat 10 times
8. Check browser memory usage (should not continuously increase)
9. Check console for any unclosed connections

### Network Throttling Test
1. Open Chrome DevTools → Network tab
2. Set throttling to "Slow 3G"
3. Navigate to dashboard
4. Verify WebSocket connects (may take longer)
5. Test real-time features still work
6. Set throttling back to "No throttling"

---

## Expected Console Output

### Successful Connection
```
🔄 Connecting to WebSocket... (Attempt 1)
✅ Dashboard WebSocket connected
```

### Receiving Messages
```
📥 Received: {type: "notification", content: "..."}
```

### Connection Closed
```
🔌 WebSocket closed: 1006 - No reason provided
⏳ Reconnecting in 1 seconds...
```

### Successful Reconnection
```
🔄 Connecting to WebSocket... (Attempt 2)
✅ Dashboard WebSocket connected
```

### Max Attempts Reached
```
❌ Max reconnection attempts reached. Please refresh the page.
```

---

## Common Issues & Solutions

### Issue: WebSocket won't connect
**Solution:**
- Verify backend server is running on port 8080
- Check backend logs for errors
- Verify no firewall blocking WebSocket connections

### Issue: Infinite reconnection loop
**Solution:**
- This should be fixed by the implementation
- If still occurring, check console for error messages
- Verify MAX_RECONNECT_ATTEMPTS is set to 10

### Issue: Messages not appearing in real-time
**Solution:**
- Check WebSocket connection status
- Verify both users are connected
- Check console for WebSocket messages
- Refresh the page

### Issue: Memory increasing over time
**Solution:**
- This should be fixed by proper cleanup
- If still occurring, check for unclosed connections
- Verify cleanup code in useEffect return functions

---

## Success Criteria

✅ All pages connect to WebSocket successfully
✅ Real-time features work on all pages
✅ Reconnection works with exponential backoff
✅ Max 10 reconnection attempts enforced
✅ Proper cleanup on page navigation
✅ No memory leaks after multiple navigations
✅ Connection status indicators work correctly
✅ No console errors during normal operation
✅ Queued messages sent after reconnection

---

## Reporting Issues

If you encounter any issues during testing, please report:
1. Page where issue occurred
2. Steps to reproduce
3. Console error messages
4. Expected vs actual behavior
5. Browser and version

---

**Last Updated:** January 12, 2025
