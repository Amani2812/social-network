# Real-Time Messaging Testing Guide

## Prerequisites
- Backend server running on http://localhost:8080
- Frontend running on http://localhost:3000
- At least 2 user accounts for testing
- Browser with Developer Console open (F12)

---

## Test 1: Normal Message Flow ✅

### Steps:
1. **Login as User A** in Browser 1
2. **Login as User B** in Browser 2 (or incognito window)
3. Navigate to Messages page on both browsers
4. User A selects User B from conversations
5. User A types and sends a message
6. **Verify**: Message appears instantly on User B's screen (< 1 second)
7. User B replies
8. **Verify**: Reply appears instantly on User A's screen

### Expected Results:
- ✅ Messages appear within 1 second
- ✅ No page refresh needed
- ✅ Messages appear in correct order
- ✅ Message input clears after sending

### Console Logs to Check:
```
User A Console:
📤 Message sent via WebSocket
📥 Messages page received: {type: "private", sender_id: A, receiver_id: B}
✅ Adding new message to conversation

User B Console:
📥 Messages page received: {type: "private", sender_id: A, receiver_id: B}
✅ Adding new message to conversation
```

---

## Test 2: Connection Status Indicator ✅

### Steps:
1. Open Messages page
2. Look at the header next to username
3. **Verify**: Green dot with "Connected" text

### Expected Results:
- ✅ Green pulsing dot visible
- ✅ Text says "Connected"
- ✅ Indicator updates in real-time

### States to Observe:
- **Connecting**: Yellow dot + "Connecting..."
- **Connected**: Green dot + "Connected" (pulsing)
- **Disconnected**: Red dot + "Disconnected" (static)

---

## Test 3: Connection Stability (5+ Minutes) ✅

### Steps:
1. Open Messages page
2. Keep page open for 5+ minutes
3. Monitor connection status indicator
4. Send a message after 5 minutes
5. Check backend server logs

### Expected Results:
- ✅ Connection stays green throughout
- ✅ Messages still send instantly after 5+ minutes
- ✅ No disconnections or reconnections

### Backend Logs to Check:
```
Every 30 seconds you should see:
Sent ping to client 123
Received pong from client 123

Should NOT see:
Client 123 disconnected
WebSocket error for client 123
```

### Console Logs to Check:
```
Should NOT see:
🔌 WebSocket closed
🔄 Reconnecting in Xms
❌ WebSocket error
```

---

## Test 4: Reconnection Logic ✅

### Steps:
1. Open Messages page (connection green)
2. **Stop the backend server** (Ctrl+C in terminal)
3. **Verify**: Connection status turns red "Disconnected"
4. Try to send a message
5. **Verify**: Message input clears (message queued)
6. **Restart the backend server**
7. **Observe**: Connection status changes: red → yellow → green
8. **Verify**: Queued message is sent automatically

### Expected Results:
- ✅ Connection status updates correctly
- ✅ Messages are queued when offline
- ✅ Automatic reconnection with exponential backoff
- ✅ Queued messages sent on reconnection

### Console Logs to Check:
```
When server stops:
🔌 WebSocket closed: 1006
🔄 Reconnecting in 1000ms (attempt 1/10)

When reconnecting:
🔄 Connecting to WebSocket...
🔄 Reconnecting in 2000ms (attempt 2/10)
🔄 Reconnecting in 4000ms (attempt 3/10)

When server restarts:
✅ WebSocket Connected
📤 Sending 1 queued messages
```

### Timing to Verify:
- Attempt 1: 1 second delay
- Attempt 2: 2 seconds delay
- Attempt 3: 4 seconds delay
- Attempt 4: 8 seconds delay
- Attempt 5: 16 seconds delay
- Attempt 6+: 30 seconds delay (max)

---

## Test 5: Message Queue (Offline Messages) ✅

### Steps:
1. Open Messages page
2. Stop backend server
3. Type and send 3 messages
4. **Verify**: All 3 messages clear from input
5. Restart backend server
6. **Verify**: All 3 messages appear in chat
7. Check other user's browser
8. **Verify**: All 3 messages received

### Expected Results:
- ✅ Messages queued when offline
- ✅ All queued messages sent on reconnection
- ✅ Messages appear in correct order
- ✅ No messages lost

### Console Logs to Check:
```
When sending offline:
⏳ WebSocket not connected, queueing message
⏳ WebSocket not connected, queueing message
⏳ WebSocket not connected, queueing message

On reconnection:
✅ WebSocket Connected
📤 Sending 3 queued messages
```

---

## Test 6: Duplicate Prevention ✅

### Steps:
1. Open Messages page
2. Send the same message twice quickly (within 2 seconds)
3. **Verify**: Both messages appear
4. Send 5 different messages quickly
5. **Verify**: All 5 messages appear

### Expected Results:
- ✅ Legitimate messages not blocked
- ✅ Only actual duplicates prevented
- ✅ 2-second window for duplicate detection

### Console Logs to Check:
```
Should see for each message:
✅ Adding new message to conversation

Should NOT see (unless actual duplicate):
⚠️ Duplicate message detected, skipping
```

---

## Test 7: Multiple Conversations ✅

### Steps:
1. Login as User A
2. Send message to User B
3. Switch to User C conversation
4. Send message to User C
5. Switch back to User B
6. **Verify**: Previous messages still visible
7. Send another message to User B

### Expected Results:
- ✅ Messages persist when switching conversations
- ✅ Each conversation maintains its own history
- ✅ Real-time updates work for all conversations

---

## Test 8: Long Connection Test ✅

### Steps:
1. Open Messages page
2. Keep page open for 10+ minutes
3. Send messages at 2-minute intervals
4. Monitor connection status
5. Check backend logs

### Expected Results:
- ✅ Connection stays green for entire duration
- ✅ All messages send instantly
- ✅ Regular ping/pong in backend logs

### Backend Logs Pattern:
```
Every 30 seconds:
Sent ping to client 123
Received pong from client 123
Sent ping to client 123
Received pong from client 123
...
```

---

## Test 9: Network Interruption Recovery ✅

### Steps:
1. Open Messages page (connected)
2. Disable network adapter or WiFi
3. **Verify**: Connection turns red
4. Try sending messages (should queue)
5. Re-enable network
6. **Verify**: Automatic reconnection
7. **Verify**: Queued messages sent

### Expected Results:
- ✅ Detects network loss
- ✅ Queues messages during outage
- ✅ Reconnects automatically
- ✅ Sends queued messages

---

## Test 10: Browser Tab Switching ✅

### Steps:
1. Open Messages page (connected)
2. Switch to another browser tab for 5 minutes
3. Switch back to Messages tab
4. **Verify**: Still connected (green)
5. Send a message
6. **Verify**: Sends instantly

### Expected Results:
- ✅ Connection maintained in background
- ✅ No reconnection needed
- ✅ Messages work immediately

---

## Common Issues and Solutions

### Issue: Connection stays yellow "Connecting..."
**Cause**: Backend server not running or wrong URL
**Solution**: 
- Check backend is running on port 8080
- Verify WebSocket URL: `ws://localhost:8080/ws`

### Issue: Connection turns red immediately
**Cause**: Authentication failure or CORS issue
**Solution**:
- Ensure user is logged in
- Check session cookie is valid
- Verify CORS settings in backend

### Issue: Messages don't appear
**Cause**: WebSocket not receiving messages
**Solution**:
- Check browser console for errors
- Verify backend logs show message broadcasting
- Check if correct conversation is selected

### Issue: Duplicate messages appearing
**Cause**: Duplicate detection too lenient
**Solution**:
- Check if messages are actually different
- Verify timestamps are correct
- Review console logs for duplicate warnings

---

## Success Criteria Checklist

After completing all tests, verify:

- [ ] Messages appear instantly (< 1 second)
- [ ] Connection stays alive for 10+ minutes
- [ ] Automatic reconnection works after network interruption
- [ ] No legitimate messages blocked as duplicates
- [ ] Connection status indicator works correctly
- [ ] Queued messages sent after reconnection
- [ ] No errors in browser console during normal operation
- [ ] Backend logs show regular ping/pong activity
- [ ] Multiple conversations work correctly
- [ ] Messages persist when switching conversations

---

## Performance Benchmarks

| Metric | Target | How to Measure |
|--------|--------|----------------|
| Message Delivery | < 1 second | Time from send to receive |
| Connection Uptime | > 99% | Monitor for 10+ minutes |
| Reconnection Time | < 5 seconds | Time from disconnect to reconnect |
| Ping Interval | 30 seconds | Check backend logs |
| Timeout Buffer | 60 seconds | 90s timeout - 30s ping |

---

## Reporting Issues

If you find any issues during testing, please note:

1. **What you were doing** (which test step)
2. **What you expected** (from Expected Results)
3. **What actually happened** (the bug)
4. **Console logs** (copy relevant logs)
5. **Backend logs** (if applicable)
6. **Screenshots** (especially of connection status)

---

## Quick Test (5 Minutes)

If you want a quick verification:

1. ✅ Open Messages page → Check green "Connected"
2. ✅ Send a message → Appears instantly
3. ✅ Wait 2 minutes → Still connected
4. ✅ Send another message → Still works
5. ✅ Stop backend → Turns red
6. ✅ Start backend → Turns green
7. ✅ Send message → Works immediately

If all 7 steps pass, the fix is working! ✅

---

**Happy Testing! 🚀**
