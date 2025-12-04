# Comprehensive Group Testing Guide

This guide covers all group-related features including invitations, join requests, member permissions, events, and notifications.

## Prerequisites

1. **Start the Backend Server**
   ```bash
   cd backend
   go run server.go
   ```

2. **Start the Frontend Server**
   ```bash
   cd frontend
   npm run dev
   ```

3. **Open Two Browser Windows**
   - Browser 1: http://localhost:3000 (User A)
   - Browser 2: http://localhost:3000 (User B) - Use Incognito/Private mode

---

## Test 1: Create Group and Invite Follower

### Setup: Follow Each Other

**Browser 1 (User A):**
1. Register/Login as User A (e.g., alice@test.com)
2. Go to Dashboard → Search for User B
3. Click on User B's profile
4. Click "Follow" button
5. ✅ **Verify**: Button changes to "Following" or "Request Pending" (depending on privacy)

**Browser 2 (User B):**
1. Register/Login as User B (e.g., bob@test.com)
2. Click on Notifications icon (🔔)
3. ✅ **Verify**: You see a follow request notification from User A
4. Click "Accept" button
5. ✅ **Verify**: Notification is marked as read
6. Go to User A's profile
7. Click "Follow" button
8. ✅ **Verify**: You are now following User A

**Browser 1 (User A):**
1. Refresh notifications
2. ✅ **Verify**: You see notification that User B is following you

### Create Group and Invite

**Browser 1 (User A):**
1. Go to Dashboard → Groups
2. Click "+ Create New Group"
3. Enter:
   - Group Name: "Test Group"
   - Description: "Testing group invitations"
4. Click "Create Group"
5. ✅ **Verify**: Group is created and you're redirected to group page
6. ✅ **Verify**: You can see the group in "Your Groups" section

**Browser 1 (User A) - Invite User B:**
1. From the group page, you need to invite User B
2. **Note**: The current UI doesn't have an "Invite" button on the group detail page
3. **Expected Feature**: There should be an "Invite Members" button that shows a list of followers

### Test Result:
- ✅ Group creation works
- ❌ **MISSING FEATURE**: No UI to invite users to groups from the group detail page
- ✅ Backend API exists: `POST /api/groups/invite`

---

## Test 2: Group Invitation Acceptance/Refusal

### Current Implementation Status:
- ✅ Backend API exists: `POST /api/groups/invite/respond`
- ✅ Notifications are created for group invites
- ❌ **MISSING**: UI to invite users from group page
- ❌ **MISSING**: UI to accept/decline group invites in notifications

### Expected Flow:
1. User A invites User B to group
2. User B receives notification
3. User B can accept or decline from notifications page
4. If accepted, User B becomes a group member

---

## Test 3: Group Join Request

**Browser 2 (User B):**
1. Go to Dashboard → Groups
2. Scroll to "Discover Groups" section
3. ✅ **Verify**: You can see "Test Group" created by User A
4. Click "View Details" button
5. ✅ **Verify**: You're redirected to the group page
6. **Note**: There's no "Request to Join" button visible

### Test Result:
- ✅ Backend API exists: `POST /api/groups/join/request`
- ❌ **MISSING FEATURE**: No "Request to Join" button on group detail page for non-members

---

## Test 4: Group Join Request Acceptance

### Current Implementation Status:
- ✅ Backend API exists: `POST /api/groups/join/respond`
- ✅ Backend API exists: `GET /api/groups/join/requests`
- ❌ **MISSING**: UI to view and respond to join requests
- ❌ **MISSING**: UI to show pending join requests to group admins

### Expected Flow:
1. User B requests to join group
2. User A (admin) receives notification
3. User A can view join requests in group settings
4. User A can accept or decline the request
5. User B receives notification of acceptance/rejection

---

## Test 5: Member Permissions - Invitations

### Test if Non-Admin Members Can Invite:

**Setup**: First, manually add User B to the group using backend API or database

**Browser 2 (User B) - As Regular Member:**
1. Go to the group page
2. **Expected**: Should see an "Invite Members" button
3. **Expected**: Should be able to invite other users

### Test Result:
- ✅ Backend allows any group member to invite: `InviteToGroup` checks if inviter is a member
- ❌ **MISSING**: UI for members to invite other users

---

## Test 6: Create Posts in Group

**Browser 1 (User A):**
1. Go to group page
2. Click on "Posts" tab
3. Type a message in the text area: "Hello from User A!"
4. Click "Post" button
5. ✅ **Verify**: Post appears in the feed
6. ✅ **Verify**: Success message: "Post created successfully! ✓"

**Browser 2 (User B) - If Member:**
1. Go to the same group page
2. ✅ **Verify**: You can see User A's post
3. Type a message: "Hello from User B!"
4. Click "Post" button
5. ✅ **Verify**: Your post appears in the feed

**Browser 2 (User B) - If NOT Member:**
1. Try to access group page
2. Try to create a post
3. ✅ **Verify**: Error message: "You might not be a member of this group"

### Test Result:
- ✅ Group posts work correctly
- ✅ Only members can post
- ❌ **MISSING**: Comment functionality on group posts

---

## Test 7: Create and Vote on Events

**Browser 1 (User A):**
1. Go to group page
2. Click on "Events" tab
3. Click "+ Create Event" button
4. ✅ **Verify**: Event creation form appears
5. Fill in:
   - Event Title: "Group Meeting"
   - Description: "Monthly sync meeting"
   - Date & Time: Select a future date/time
6. Click "Create Event"
7. ✅ **Verify**: Event is created and appears in the list
8. ✅ **Verify**: Event shows title, description, and date/time

**Browser 2 (User B) - If Member:**
1. Go to the same group page
2. Click on "Events" tab
3. ✅ **Verify**: You can see the "Group Meeting" event
4. ✅ **Verify**: You see two buttons: "✓ Going" and "✗ Not Going"
5. Click "✓ Going" button
6. ✅ **Verify**: Alert shows "You responded: Going ✓"
7. Click "✗ Not Going" button
8. ✅ **Verify**: Alert shows "You responded: Not Going ✗"

### Test Result:
- ✅ Event creation works
- ✅ Event voting works
- ✅ Members can see events
- ❌ **MISSING**: Display of who is going/not going
- ❌ **MISSING**: Event response count

---

## Test 8: Event Notifications

**Browser 2 (User B) - After Event Creation:**
1. Click on Notifications icon (🔔)
2. ✅ **Verify**: You should see a notification about the new event
3. ✅ **Verify**: Notification content: "New event in Test Group: Group Meeting"
4. ✅ **Verify**: Notification icon: 📅
5. Click on the notification
6. ❌ **ISSUE**: Navigation to `/events/{id}` doesn't exist (should go to group page)

### Test Result:
- ✅ Event notifications are created
- ✅ Notifications appear in the notifications page
- ❌ **BUG**: Clicking event notification tries to navigate to non-existent `/events/{id}` page

---

## Test 9: Group Invitation Notifications

### Current Status:
- ✅ Backend creates notifications: `CreateNotification(userID, "group_invite", content, &groupID)`
- ✅ Notifications appear in notifications page
- ❌ **MISSING**: Accept/Decline buttons for group invites in notifications
- ✅ Clicking notification navigates to group page

### Expected Behavior:
1. User receives group invitation
2. Notification shows: "You've been invited to join {Group Name}"
3. Notification has Accept/Decline buttons
4. Clicking Accept adds user to group
5. Clicking Decline removes the invitation

---

## Test 10: Group Join Request Notifications

### Current Status:
- ✅ Backend creates notifications for group admins
- ✅ Notification content: "{User} wants to join {Group}"
- ❌ **MISSING**: UI to display join requests to admins
- ❌ **MISSING**: Accept/Decline functionality in notifications

### Expected Behavior:
1. User B requests to join group
2. User A (admin) receives notification
3. Notification shows: "Bob Smith wants to join Test Group"
4. Notification has Accept/Decline buttons
5. Clicking Accept adds User B to group
6. User B receives acceptance notification

---

## Summary of Test Results

### ✅ Working Features:
1. ✅ Group creation
2. ✅ Viewing all groups
3. ✅ Group posts (create and view)
4. ✅ Event creation
5. ✅ Event voting (Going/Not Going)
6. ✅ Event notifications are created
7. ✅ Follow request notifications with Accept/Decline buttons
8. ✅ Member-only access to group posts
9. ✅ Backend APIs for all group features

### ❌ Missing Features/Bugs:

#### Critical Missing Features:
1. ❌ **No UI to invite users to groups** from group detail page
2. ❌ **No "Request to Join" button** for non-members on group page
3. ❌ **No UI to view/respond to group join requests** for admins
4. ❌ **No Accept/Decline buttons for group invites** in notifications
5. ❌ **No comment functionality** on group posts

#### Bugs:
1. ❌ **Event notification navigation bug**: Clicking event notification tries to go to `/events/{id}` which doesn't exist (should go to group page with events tab)

#### Missing Enhancements:
1. ❌ **No display of event responses** (who's going/not going)
2. ❌ **No member list** visible on group page
3. ❌ **No group settings/management** page for admins
4. ❌ **No way to leave a group**
5. ❌ **No way to remove members** (for admins)

---

## Recommended Fixes

### Priority 1 (Critical for Basic Functionality):

1. **Add "Invite Members" Button to Group Page**
   - Location: Group detail page header
   - Shows list of followers
   - Allows selecting and inviting users

2. **Add "Request to Join" Button**
   - Location: Group detail page for non-members
   - Visible in "Discover Groups" section
   - Sends join request to group admins

3. **Add Group Invite Accept/Decline in Notifications**
   - Add buttons similar to follow requests
   - Handle `group_invite` notification type
   - Call `/api/groups/invite/respond` API

4. **Add Join Request Management UI**
   - Add "Pending Requests" section for admins
   - Show list of users requesting to join
   - Accept/Decline buttons for each request

5. **Fix Event Notification Navigation**
   - Change navigation from `/events/{id}` to `/groups/{groupId}?tab=events`
   - Or create a dedicated events page

### Priority 2 (Enhanced Functionality):

1. **Add Comments to Group Posts**
   - Similar to regular posts
   - Show comment count
   - Allow members to comment

2. **Add Member List**
   - Show all group members
   - Display member roles (admin/member)
   - Show member count

3. **Add Event Response Display**
   - Show count of "Going" and "Not Going"
   - Show list of users and their responses
   - Highlight current user's response

4. **Add Group Management Features**
   - Leave group button
   - Remove member button (admin only)
   - Edit group details (admin only)
   - Delete group (creator only)

---

## Testing Checklist

Use this checklist to verify all features:

### Group Creation & Discovery
- [ ] Can create a new group
- [ ] Group appears in "Your Groups"
- [ ] Group appears in "Discover Groups" for other users
- [ ] Can view group details

### Group Invitations
- [ ] Can invite followers to group
- [ ] Invited user receives notification
- [ ] Invited user can accept invitation
- [ ] Invited user can decline invitation
- [ ] Accepted user becomes group member
- [ ] Declined invitation is removed

### Group Join Requests
- [ ] Non-member can request to join
- [ ] Admin receives join request notification
- [ ] Admin can view pending requests
- [ ] Admin can accept join request
- [ ] Admin can decline join request
- [ ] Requester receives acceptance/rejection notification

### Member Permissions
- [ ] Regular members can invite users
- [ ] Regular members can create posts
- [ ] Regular members can create events
- [ ] Regular members can comment on posts
- [ ] Non-members cannot post
- [ ] Non-members cannot see group content

### Group Posts
- [ ] Members can create posts
- [ ] Posts appear in chronological order
- [ ] Posts show author information
- [ ] Members can comment on posts
- [ ] Post images are displayed correctly

### Group Events
- [ ] Members can create events
- [ ] Events show title, description, date/time
- [ ] Members can vote "Going"
- [ ] Members can vote "Not Going"
- [ ] Event responses are saved
- [ ] Event response counts are displayed
- [ ] List of attendees is visible

### Notifications
- [ ] Group invite notifications appear
- [ ] Group invite notifications have Accept/Decline buttons
- [ ] Join request notifications appear for admins
- [ ] Join request notifications have Accept/Decline buttons
- [ ] Event creation notifications appear for members
- [ ] Clicking notifications navigates correctly
- [ ] Notifications are marked as read

---

## API Endpoints Reference

### Group Management
- `POST /api/groups/create` - Create a new group
- `GET /api/groups/get?id={id}` - Get group details
- `GET /api/groups/user` - Get user's groups
- `GET /api/groups/all` - Get all groups

### Group Membership
- `POST /api/groups/invite` - Invite user to group
- `POST /api/groups/invite/respond` - Accept/decline group invite
- `POST /api/groups/join/request` - Request to join group
- `GET /api/groups/join/requests?group_id={id}` - Get join requests (admin)
- `POST /api/groups/join/respond` - Accept/decline join request (admin)

### Group Content
- `POST /api/groups/posts/create` - Create group post
- `GET /api/groups/posts/get?group_id={id}` - Get group posts
- `POST /api/events/create` - Create event
- `GET /api/events/get?group_id={id}` - Get group events
- `POST /api/events/respond` - Respond to event (going/not_going)

### Notifications
- `GET /api/notifications` - Get user notifications
- `POST /api/notifications/read` - Mark notification as read
- `GET /api/notifications/unread` - Get unread count

---

## Conclusion

The social network application has a solid foundation for group functionality with most backend APIs implemented. However, several critical UI components are missing that prevent users from fully utilizing the group features. The main gaps are:

1. No UI for inviting users to groups
2. No UI for requesting to join groups
3. No UI for managing join requests
4. No accept/decline buttons for group invitations in notifications

Once these UI components are added, the group functionality will be complete and fully testable according to the requirements.
