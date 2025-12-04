# Comprehensive Testing Guide - Social Network Features

This guide provides step-by-step instructions for testing all features of the social network application.

## Prerequisites

1. **Two or Three Browsers**: Use different browsers (Chrome, Firefox, Edge) or incognito/private windows
2. **Test Users**: Create at least 3 test users with different profiles
3. **Backend Running**: Ensure the Go backend is running on `http://localhost:8080`
4. **Frontend Running**: Ensure the Next.js frontend is running on `http://localhost:3000`

## Test User Setup

Create these test users for comprehensive testing:

1. **User A (Public Profile)**
   - Email: alice@test.com
   - Name: Alice Johnson
   - Profile: Public

2. **User B (Private Profile)**
   - Email: bob@test.com
   - Name: Bob Smith
   - Profile: Private

3. **User C (Public Profile)**
   - Email: charlie@test.com
   - Name: Charlie Brown
   - Profile: Public

---

## 1. Follow System Tests

### Test 1.1: Follow Public User
**Objective**: Verify that following a public user works instantly

**Steps**:
1. Open Browser 1, log in as User A
2. Open Browser 2, log in as User B
3. In Browser 1 (User A), search for User B
4. Click "Follow" button on User B's profile
5. In Browser 2 (User B), check notifications

**Expected Results**:
- ✅ User A should immediately follow User B (no approval needed)
- ✅ User B should receive a notification: "Alice Johnson started following you"
- ✅ User B should see User A in their followers list

### Test 1.2: Follow Private User
**Objective**: Verify that following a private user requires approval

**Steps**:
1. In Browser 1 (User A), search for User C (ensure User C has private profile)
2. Click "Follow" button on User C's profile
3. In Browser 3 (User C), check notifications
4. User C clicks on the notification
5. User C accepts the follow request

**Expected Results**:
- ✅ Follow button should show "Pending" after clicking
- ✅ User C should receive notification: "Alice Johnson wants to follow you"
- ✅ User C should see accept/reject buttons
- ✅ After accepting, User A should appear in User C's followers
- ✅ User A should see User C's "almost private" posts

### Test 1.3: Reject Follow Request
**Objective**: Verify that rejecting a follow request works

**Steps**:
1. User A sends follow request to private User C
2. User C receives notification
3. User C clicks "Reject"

**Expected Results**:
- ✅ User A should not appear in User C's followers
- ✅ User A should not see User C's private posts
- ✅ Follow button should return to "Follow" state

---

## 2. Group Management Tests

### Test 2.1: Create Group
**Objective**: Verify group creation

**Steps**:
1. Log in as User A
2. Navigate to Groups page
3. Click "Create Group"
4. Enter group name: "Test Group"
5. Enter description: "A test group"
6. Submit

**Expected Results**:
- ✅ Group should be created successfully
- ✅ User A should be the admin
- ✅ Group should appear in User A's groups list

### Test 2.2: Invite User to Group (Admin)
**Objective**: Verify that admins can invite users

**Steps**:
1. User A (admin) opens the group
2. Clicks "Invite Members"
3. Searches for User B
4. Sends invitation
5. User B checks notifications

**Expected Results**:
- ✅ User B receives notification: "You've been invited to join Test Group"
- ✅ User B can see accept/reject options
- ✅ Invitation appears as "pending" in group members

### Test 2.3: Accept Group Invitation
**Objective**: Verify accepting group invitations

**Steps**:
1. User B clicks on group invitation notification
2. Clicks "Accept"

**Expected Results**:
- ✅ User B becomes a member of the group
- ✅ User B can see group posts
- ✅ User B appears in members list with role "member"

### Test 2.4: Request to Join Group
**Objective**: Verify users can request to join groups

**Steps**:
1. User C (not a member) navigates to "Test Group"
2. Clicks "Request to Join" button
3. User A (admin) checks notifications

**Expected Results**:
- ✅ Button changes to "Request Pending"
- ✅ User A receives notification: "Charlie Brown wants to join Test Group"
- ✅ User A can see the request in "Pending Requests" section

### Test 2.5: Accept Join Request
**Objective**: Verify admins can accept join requests

**Steps**:
1. User A opens group page
2. Views "Pending Requests" section
3. Clicks "Accept" on User C's request
4. User C checks notifications

**Expected Results**:
- ✅ User C receives notification: "Your request to join Test Group has been accepted"
- ✅ User C becomes a member
- ✅ User C can now see and create posts in the group

### Test 2.6: Reject Join Request
**Objective**: Verify admins can reject join requests

**Steps**:
1. User D requests to join group
2. User A (admin) clicks "Reject"
3. User D checks notifications

**Expected Results**:
- ✅ User D receives notification: "Your request to join Test Group has been declined"
- ✅ User D is not added to the group
- ✅ User D can request again if desired

### Test 2.7: Member Invites Another User
**Objective**: Verify that regular members can invite users

**Steps**:
1. User B (regular member, not admin) opens group
2. Clicks "Invite Members"
3. Invites User D
4. User D checks notifications

**Expected Results**:
- ✅ User B can access invite functionality
- ✅ User D receives invitation notification
- ✅ Invitation works same as admin invitations

### Test 2.8: Create Group Post
**Objective**: Verify members can create posts

**Steps**:
1. User B (member) opens group
2. Creates a post: "Hello from User B!"
3. User A and User C check the group

**Expected Results**:
- ✅ Post appears in group feed
- ✅ All members can see the post
- ✅ Post shows User B as author

### Test 2.9: Comment on Group Post
**Objective**: Verify members can comment on group posts

**Steps**:
1. User C opens the group
2. Finds User B's post
3. Adds comment: "Great post!"

**Expected Results**:
- ✅ Comment appears under the post
- ✅ All members can see the comment
- ✅ Comment shows User C as author

---

## 3. Real-Time Messaging Tests

### Test 3.1: Private Message Between Two Users
**Objective**: Verify real-time private messaging

**Steps**:
1. Browser 1: User A opens Messages page
2. Browser 2: User B opens Messages page
3. User A selects User B from conversations
4. User A types: "Hello Bob!" and sends
5. Observe Browser 2 (User B)

**Expected Results**:
- ✅ Message appears instantly in User B's chat (real-time)
- ✅ Message shows correct sender (User A)
- ✅ Message appears in correct conversation
- ✅ No server crash

### Test 3.2: Continuous Chat Between Two Users
**Objective**: Verify chat stability

**Steps**:
1. User A and User B exchange 10+ messages
2. Include emojis: 😊 👍 ❤️
3. Send messages rapidly

**Expected Results**:
- ✅ All messages delivered in correct order
- ✅ Emojis display correctly
- ✅ No messages lost
- ✅ Server remains stable
- ✅ No crashes

### Test 3.3: Private Message to Specific User (Three Users)
**Objective**: Verify messages only go to intended recipient

**Steps**:
1. Browser 1: User A logged in
2. Browser 2: User B logged in
3. Browser 3: User C logged in
4. User A sends message to User B: "Message for Bob only"
5. Check all three browsers

**Expected Results**:
- ✅ Only User B receives the message
- ✅ User C does NOT see the message
- ✅ Message appears in User A's chat (echo)

### Test 3.4: Group Chat - Real-Time
**Objective**: Verify real-time group messaging

**Steps**:
1. Browser 1: User A in "Test Group"
2. Browser 2: User B in "Test Group"
3. Browser 3: User C in "Test Group"
4. User A sends message: "Hello everyone!"
5. Observe all browsers

**Expected Results**:
- ✅ All group members receive message instantly
- ✅ Message shows correct sender
- ✅ Real-time delivery (no refresh needed)

### Test 3.5: Group Chat Stability
**Objective**: Verify group chat doesn't crash

**Steps**:
1. All three users in same group chat
2. Send 20+ messages rapidly
3. Include emojis and longer messages
4. Multiple users send simultaneously

**Expected Results**:
- ✅ All messages delivered
- ✅ Correct order maintained
- ✅ No server crashes
- ✅ All users can continue chatting

### Test 3.6: Emoji Support
**Objective**: Verify emoji support in messages

**Steps**:
1. User A sends to User B: "Hello 😊 How are you? 👍"
2. User B replies: "Great! ❤️🎉"
3. Send emojis in group chat

**Expected Results**:
- ✅ Emojis display correctly in private chat
- ✅ Emojis display correctly in group chat
- ✅ Emojis saved and retrieved correctly

---

## 4. Notification System Tests

### Test 4.1: Follow Request Notification
**Objective**: Verify follow request notifications

**Steps**:
1. User A follows private User C
2. User C checks notifications page

**Expected Results**:
- ✅ Notification appears: "Alice Johnson wants to follow you"
- ✅ Notification shows as unread
- ✅ Clicking notification navigates to appropriate page
- ✅ Accept/reject buttons work

### Test 4.2: Group Invitation Notification
**Objective**: Verify group invitation notifications

**Steps**:
1. User A invites User B to group
2. User B checks notifications

**Expected Results**:
- ✅ Notification appears: "You've been invited to join [Group Name]"
- ✅ Notification is unread
- ✅ Clicking navigates to group page
- ✅ Accept/reject options available

### Test 4.3: Group Join Request Notification
**Objective**: Verify join request notifications for admins

**Steps**:
1. User C requests to join group
2. User A (admin) checks notifications

**Expected Results**:
- ✅ Notification appears: "Charlie Brown wants to join [Group Name]"
- ✅ Notification is unread
- ✅ Clicking navigates to group with pending requests visible
- ✅ Accept/reject buttons work

### Test 4.4: Event Creation Notification
**Objective**: Verify event notifications

**Steps**:
1. User A creates event in group
2. All group members check notifications

**Expected Results**:
- ✅ All members receive notification: "New event in [Group]: [Event Name]"
- ✅ Notifications are unread
- ✅ Clicking navigates to event details

### Test 4.5: Notification on Every Page
**Objective**: Verify notifications accessible from all pages

**Steps**:
1. Generate some notifications
2. Navigate to different pages:
   - Dashboard
   - Profile
   - Groups
   - Messages
   - Search

**Expected Results**:
- ✅ Notification bell/icon visible on all pages
- ✅ Unread count displays correctly
- ✅ Clicking opens notifications
- ✅ Notifications accessible from anywhere

---

## 5. Integration Tests

### Test 5.1: Complete User Journey
**Objective**: Test complete workflow

**Steps**:
1. User A creates account
2. User A creates a group
3. User A invites User B
4. User B accepts invitation
5. User C requests to join
6. User A accepts request
7. All users post in group
8. All users comment on posts
9. User A creates event
10. All users receive notifications
11. Users chat in group
12. Users send private messages

**Expected Results**:
- ✅ All steps complete successfully
- ✅ No errors or crashes
- ✅ All notifications delivered
- ✅ All messages delivered in real-time

### Test 5.2: Concurrent Operations
**Objective**: Test system under concurrent load

**Steps**:
1. Three users simultaneously:
   - Send messages
   - Create posts
   - Send notifications
   - Accept/reject requests

**Expected Results**:
- ✅ All operations complete successfully
- ✅ No race conditions
- ✅ No data corruption
- ✅ Server remains stable

---

## 6. Edge Cases and Error Handling

### Test 6.1: Duplicate Requests
**Steps**:
1. User A sends follow request to User C
2. User A tries to send another follow request

**Expected Results**:
- ✅ System prevents duplicate requests
- ✅ Appropriate error message shown

### Test 6.2: Invalid Operations
**Steps**:
1. Non-member tries to post in group
2. Non-admin tries to accept join requests
3. User tries to invite themselves

**Expected Results**:
- ✅ Operations blocked with appropriate errors
- ✅ No system crashes

### Test 6.3: Offline User Messages
**Steps**:
1. User B logs out
2. User A sends message to User B
3. User B logs back in

**Expected Results**:
- ✅ Message saved in database
- ✅ User B sees message when logging back in
- ✅ Notification created for offline user

---

## Testing Checklist Summary

### Follow System
- [ ] Follow public user instantly
- [ ] Follow private user with approval
- [ ] Reject follow request
- [ ] Unfollow user
- [ ] View followers list
- [ ] View following list

### Group Management
- [ ] Create group
- [ ] Admin invites user
- [ ] Member invites user
- [ ] User requests to join
- [ ] Admin accepts join request
- [ ] Admin rejects join request
- [ ] Accept group invitation
- [ ] Reject group invitation
- [ ] Create group post
- [ ] Comment on group post

### Real-Time Messaging
- [ ] Private message delivery (real-time)
- [ ] Private message to correct user only
- [ ] Group message delivery (real-time)
- [ ] Chat stability (no crashes)
- [ ] Emoji support
- [ ] Multiple concurrent chats

### Notifications
- [ ] Follow request notification
- [ ] Group invitation notification
- [ ] Join request notification
- [ ] Event creation notification
- [ ] Notifications on all pages
- [ ] Mark as read functionality
- [ ] Unread count accuracy

### Integration
- [ ] Complete user journey
- [ ] Concurrent operations
- [ ] Edge cases handled
- [ ] Error messages appropriate

---

## Reporting Issues

When reporting issues, include:
1. **Test Case**: Which test failed
2. **Steps**: Exact steps to reproduce
3. **Expected**: What should happen
4. **Actual**: What actually happened
5. **Browser**: Which browser(s)
6. **Screenshots**: If applicable
7. **Console Errors**: Any errors in browser console or server logs

---

## Success Criteria

All tests should pass with:
- ✅ No server crashes
- ✅ No data loss
- ✅ Real-time delivery working
- ✅ All notifications delivered
- ✅ Proper error handling
- ✅ Good user experience
