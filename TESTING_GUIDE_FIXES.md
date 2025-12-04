# Testing Guide for Member Management and WebSocket Fixes

## Prerequisites
- Frontend server running on http://localhost:3001
- Backend server running on http://localhost:8080
- At least 2 user accounts for testing
- At least 1 group created by one of the users

## Test 1: Member Management Tab Visibility

### Objective
Verify that the "Manage Members" tab appears consistently for group creators/admins across all browsers.

### Steps:

#### Chrome Testing
1. Open Chrome browser
2. Navigate to http://localhost:3001
3. Login with a user account that created a group
4. Navigate to Groups page
5. Click on a group you created
6. **Expected Result**: You should see three tabs: "📝 Posts", "📅 Events", and "👥 Manage Members"
7. Click on "Manage Members" tab
8. **Expected Result**: Tab should display join requests (if any) or "No pending join requests"

#### Firefox Testing
1. Open Firefox browser
2. Repeat steps 2-8 from Chrome testing
3. **Expected Result**: Same behavior as Chrome

#### Edge Testing
1. Open Edge browser
2. Repeat steps 2-8 from Chrome testing
3. **Expected Result**: Same behavior as Chrome

### Success Criteria
- ✅ "Manage Members" tab is visible in all three browsers
- ✅ Tab is clickable and displays content
- ✅ No errors in browser console related to member management

---

## Test 2: WebSocket Console Errors

### Objective
Verify that WebSocket connection errors no longer appear as RED errors in the console.

### Steps:

#### Dashboard Page
1. Open Chrome browser (or any browser)
2. Press F12 to open Developer Tools
3. Go to the "Console" tab
4. Navigate to http://localhost:3001
5. Login to your account
6. You should be on the Dashboard page
7. **Check Console**: Look for WebSocket-related messages
8. **Expected Result**: 
   - ✅ Should see: `⚠️ WebSocket connection error (this is normal during reconnection):` (in regular log color, NOT red)
   - ❌ Should NOT see: `❌ WebSocket error:` (in red error color)

#### Notifications Page
1. With Developer Tools still open
2. Click on the Notifications bell icon or navigate to /notifications
3. **Check Console**: Look for WebSocket-related messages
4. **Expected Result**: 
   - ✅ Should see: `⚠️ WebSocket connection error (this is normal during reconnection):` (in regular log color)
   - ❌ Should NOT see red error messages

#### Messages Page
1. With Developer Tools still open
2. Navigate to Messages page (/messages)
3. **Check Console**: Look for WebSocket-related messages
4. **Expected Result**: 
   - ✅ Should see: `⚠️ WebSocket connection error (this is normal during reconnection):` (in regular log color)
   - ❌ Should NOT see red error messages

### Success Criteria
- ✅ No RED WebSocket errors in console on any page
- ✅ Informative log messages appear instead (with ⚠️ warning icon)
- ✅ WebSocket connections still work (real-time features functional)

---

## Test 3: Verify Real-Time Features Still Work

### Objective
Ensure that fixing the console errors didn't break real-time functionality.

### Steps:

#### Real-Time Notifications
1. Open two browser windows side by side
2. Login as User A in Window 1
3. Login as User B in Window 2
4. In Window 2 (User B), navigate to User A's profile
5. Click "Follow" button
6. **Expected Result**: User A (Window 1) should receive a real-time notification

#### Real-Time Messaging
1. Keep both browser windows open
2. In Window 1 (User A), navigate to Messages
3. In Window 2 (User B), navigate to Messages
4. User B sends a message to User A
5. **Expected Result**: User A should see the message appear in real-time without refreshing

### Success Criteria
- ✅ Notifications appear in real-time
- ✅ Messages appear in real-time
- ✅ No errors in console during real-time operations

---

## Test 4: Edge Cases for Member Management

### Objective
Test the fallback mechanism for admin detection.

### Steps:

#### Test Network Failure Scenario
1. Open Chrome Developer Tools (F12)
2. Go to "Network" tab
3. Enable "Offline" mode (throttling dropdown)
4. Navigate to a group you created
5. Disable "Offline" mode
6. Refresh the page
7. **Expected Result**: "Manage Members" tab should still appear (fallback to creator check)

#### Test Non-Admin User
1. Login as a user who did NOT create any groups
2. Navigate to Groups page
3. Click on a group you're a member of (but didn't create)
4. **Expected Result**: Should NOT see "Manage Members" tab (only "Posts" and "Events")

### Success Criteria
- ✅ Group creators always see "Manage Members" tab
- ✅ Non-admin members don't see "Manage Members" tab
- ✅ Fallback mechanism works when API fails

---

## Reporting Results

### If All Tests Pass
- Document: "All fixes verified successfully"
- Note any observations or improvements

### If Any Test Fails
- Document which test failed
- Provide screenshots of console errors
- Describe the unexpected behavior
- Note the browser and version where it failed

---

## Quick Verification Checklist

- [ ] Member Management tab visible in Chrome
- [ ] Member Management tab visible in Firefox
- [ ] Member Management tab visible in Edge
- [ ] No red WebSocket errors on Dashboard
- [ ] No red WebSocket errors on Notifications page
- [ ] No red WebSocket errors on Messages page
- [ ] Real-time notifications work
- [ ] Real-time messaging works
- [ ] Non-admin users don't see Manage Members tab
- [ ] Fallback mechanism works for group creators

---

## Notes

- The frontend is running on port 3001 (not 3000)
- WebSocket connection attempts during page load are normal
- The warning messages in console are informative, not errors
- All real-time features should continue working normally
