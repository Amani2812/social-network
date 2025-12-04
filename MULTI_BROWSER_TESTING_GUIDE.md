# Multi-Browser Real-Time Notifications Testing Guide

## Overview
This guide will help you test the real-time notification system for group join requests and event creation using two different browsers.

## Prerequisites
- Backend server running on `http://localhost:8080`
- Frontend running on `http://localhost:3000`
- Two different browsers (e.g., Chrome and Firefox) or two incognito/private windows

## Test Setup

### Step 1: Prepare Two Browsers
1. **Browser A (Chrome)** - User A (Group Admin/Event Creator)
2. **Browser B (Firefox/Edge)** - User B (Requester/Member)

### Step 2: Create Test Users
If you don't have test users, create two accounts:

**User A:**
- Email: `usera@test.com`
- Password: `password123`
- Name: Alice Admin

**User B:**
- Email: `userb@test.com`
- Password: `password123`
- Name: Bob Member

## Test Scenario 1: Group Join Request Notifications

### Objective
Verify that when User B requests to join a group created by User A, User A receives a real-time notification.

### Steps:

#### Browser A (User A - Alice):
1. Open Chrome and navigate to `http://localhost:3000`
2. Login as User A (`usera@test.com`)
3. Navigate to **Groups** page
4. Click **"Create Group"**
5. Fill in:
   - Title: "Test Group"
   - Description: "Testing real-time notifications"
6. Click **"Create"**
7. **Keep this browser open and visible** - watch for notifications

#### Browser B (User B - Bob):
1. Open Firefox/Edge and navigate to `http://localhost:3000`
2. Login as User B (`userb@test.com`)
3. Navigate to **Groups** page
4. Find "Test Group" in the list
5. Click **"Request to Join"**
6. Confirm the request was sent

#### Expected Results in Browser A:
✅ **Real-time notification should appear** in the top-right corner:
   - Toast notification: "Alice Admin wants to join Test Group"
   - Notification bell badge count increases
   - No page refresh required

✅ **In Notifications Page:**
   - New notification appears at the top
   - Marked as unread (blue background)
   - Shows timestamp

### Verification Checklist:
- [ ] Toast notification appeared in Browser A without refresh
- [ ] Notification bell badge updated in real-time
- [ ] Notification visible in notifications page
- [ ] Notification contains correct user name and group name
- [ ] Notification is marked as unread

---

## Test Scenario 2: Event Creation Notifications

### Objective
Verify that when User A creates an event in a group where both users are members, User B receives a real-time notification.

### Prerequisites:
- Both User A and User B must be members of the same group
- If using the group from Test 1, User A should accept User B's join request first

### Steps:

#### Browser A (User A - Accept Join Request):
1. In Browser A, click the notification bell
2. Find Bob's join request notification
3. Click **"Accept"** (if using Test 1 group)
4. Confirm Bob is now a member

#### Browser B (User B - Verify Membership):
1. In Browser B, refresh the Groups page
2. Verify "Test Group" now shows you as a member
3. Click on "Test Group" to enter
4. **Keep this browser open and visible** - watch for notifications

#### Browser A (User A - Create Event):
1. In Browser A, navigate to "Test Group"
2. Click on **"Events"** tab
3. Click **"+ Create Event"**
4. Fill in:
   - Title: "Team Meeting"
   - Description: "Testing event notifications"
   - Date & Time: Select any future date/time
5. Click **"Create Event"**

#### Expected Results in Browser B:
✅ **Real-time notification should appear** immediately:
   - Toast notification: "New event in Test Group: Team Meeting"
   - Notification bell badge increases
   - No page refresh required

✅ **In Group Events Tab:**
   - Event appears in the list (may need to refresh to see details)

✅ **In Notifications Page:**
   - New notification for event creation
   - Marked as unread
   - Clickable to navigate to event

### Verification Checklist:
- [ ] Toast notification appeared in Browser B without refresh
- [ ] Notification bell badge updated in real-time
- [ ] Event notification visible in notifications page
- [ ] Notification contains correct group name and event title
- [ ] Clicking notification navigates to correct location

---

## Test Scenario 3: Multiple Notifications

### Objective
Test that multiple notifications work correctly in sequence.

### Steps:
1. **Browser B**: Request to join another group (or create a new group and have Browser A request to join)
2. **Browser A**: Create another event in the same group
3. **Verify**: Both notifications appear in real-time in the respective browsers

### Expected Results:
✅ All notifications appear in real-time
✅ Notification count updates correctly
✅ Notifications stack properly in the toast area
✅ Each notification auto-dismisses after 5 seconds

---

## Troubleshooting

### Issue: No Real-Time Notifications Appearing

**Check 1: WebSocket Connection**
1. Open Browser Developer Tools (F12)
2. Go to Console tab
3. Look for: `✅ Dashboard WebSocket connected` or `✅ Notifications WebSocket connected`
4. If you see connection errors, check if backend is running

**Check 2: Backend Logs**
1. Check backend terminal for WebSocket connections
2. Look for: `Client registered: <user_id>`
3. Verify notifications are being sent

**Check 3: Network Tab**
1. Open Developer Tools → Network tab
2. Filter by "WS" (WebSocket)
3. Click on the WebSocket connection
4. Check "Messages" tab to see real-time data flow

### Issue: Notifications Appear Only After Refresh

**Possible Causes:**
- WebSocket not connected
- Browser blocking WebSocket connections
- Firewall/antivirus blocking connections

**Solutions:**
1. Check browser console for errors
2. Try disabling browser extensions
3. Check if localhost is allowed in firewall

### Issue: Notification Count Not Updating

**Check:**
1. Verify `fetchUnreadCount()` is being called in WebSocket handler
2. Check browser console for API errors
3. Verify backend `/api/notifications/unread` endpoint is working

---

## Success Criteria

### ✅ Test Passed If:
1. **Group Join Request:**
   - Admin receives real-time notification when someone requests to join
   - No page refresh required
   - Notification appears within 1-2 seconds

2. **Event Creation:**
   - All group members receive real-time notification
   - No page refresh required
   - Notification appears within 1-2 seconds

3. **Notification System:**
   - Toast notifications appear and auto-dismiss
   - Bell badge updates in real-time
   - Notifications persist in database
   - Notifications page shows all notifications
   - Clicking notifications navigates correctly

### ❌ Test Failed If:
- Notifications only appear after page refresh
- WebSocket connection fails
- Notifications don't appear at all
- Notification count doesn't update
- Toast notifications don't appear

---

## Additional Testing

### Test Real-Time Updates on Notifications Page
1. Open notifications page in Browser A
2. Have Browser B send a join request
3. **Expected**: Notification list updates automatically without refresh

### Test Multiple Group Members
1. Create a group with 3+ members
2. Have one member create an event
3. **Expected**: All other members receive notifications simultaneously

### Test Notification Persistence
1. Receive a notification
2. Close browser
3. Reopen and login
4. **Expected**: Notification still visible in notifications page

---

## Demo Script

For a quick demonstration:

```bash
# Terminal 1: Start Backend
cd backend
go run server.go

# Terminal 2: Start Frontend
cd frontend
npm run dev

# Browser A (Chrome):
# 1. Login as usera@test.com
# 2. Create group "Demo Group"
# 3. Watch for notifications

# Browser B (Firefox):
# 1. Login as userb@test.com
# 2. Request to join "Demo Group"
# 3. Watch Browser A for real-time notification!

# Browser A:
# 1. Accept join request
# 2. Create event "Demo Event"
# 3. Watch Browser B for real-time notification!
```

---

## Notes

- **WebSocket Reconnection**: If connection drops, it automatically reconnects after 3 seconds
- **Toast Duration**: Notifications auto-dismiss after 5 seconds
- **Notification Types**: System supports multiple types (follow_request, group_join_request, event_invite, etc.)
- **Real-Time Updates**: Works across all pages with WebSocket connections (Dashboard, Notifications, Groups)

---

## Conclusion

This testing guide ensures that the real-time notification system works correctly for group join requests and event creation. Follow the steps carefully and verify each checkpoint to ensure full functionality.

**Happy Testing! 🎉**
