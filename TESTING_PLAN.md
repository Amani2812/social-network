# Social Network Application - Comprehensive Testing Plan

## Prerequisites
- Ensure the application is running (backend on port 8080, frontend on port 3000)
- Have at least 2 different browsers available (Chrome, Firefox, Edge, etc.)
- Clear browser cookies/cache before starting tests

## Test Environment Setup

### Starting the Application
```bash
# Option 1: Using Docker
docker-compose up --build

# Option 2: Manual start
# Terminal 1 - Backend
cd backend
go run server.go

# Terminal 2 - Frontend
cd frontend
npm install
npm run dev
```

---

## 1. FOLLOW/UNFOLLOW FUNCTIONALITY TESTS

### Test 1.1: Follow Request to Private User
**Objective**: Verify that following a private user requires sending a follow request

**Steps**:
1. Open Browser 1 (Chrome) and register/login as User A (private profile)
2. Open Browser 2 (Firefox) and register/login as User B
3. With User B, search for User A and navigate to their profile
4. Click the "Follow" button
5. Verify that the button changes to "Request Pending"
6. Switch to Browser 1 (User A)
7. Navigate to notifications or pending requests section
8. Verify that User A sees a follow request from User B

**Expected Results**:
- ✅ Follow request is sent (not auto-followed)
- ✅ User B sees "Request Pending" status
- ✅ User A receives the follow request notification
- ✅ User A can accept or decline the request

---

### Test 1.2: Follow Public User (Auto-Follow)
**Objective**: Verify that following a public user happens immediately without request

**Steps**:
1. Open Browser 1 and login as User A (ensure profile is PUBLIC)
2. Open Browser 2 and login as User B
3. With User B, navigate to User A's profile
4. Click the "Follow" button
5. Verify the button immediately changes to "Unfollow"
6. Check User A's followers list

**Expected Results**:
- ✅ User B immediately follows User A (no pending request)
- ✅ Button changes to "Unfollow" instantly
- ✅ User B appears in User A's followers list
- ✅ User A's posts appear in User B's feed

---

### Test 1.3: Accept/Decline Follow Request (Private User)
**Objective**: Verify that private users can accept or decline follow requests

**Steps**:
1. Browser 1: Login as User A (private profile)
2. Browser 2: Login as User B
3. User B sends follow request to User A (see Test 1.1)
4. User A navigates to pending requests/notifications
5. User A clicks "Accept" on User B's request
6. Verify User B now follows User A
7. Repeat with User C, but this time User A clicks "Decline"

**Expected Results**:
- ✅ User A sees pending follow requests
- ✅ User A can accept requests (User B becomes follower)
- ✅ User A can decline requests (User C does not become follower)
- ✅ Accepted followers can see User A's posts
- ✅ Declined users cannot see private content

---

### Test 1.4: Unfollow User
**Objective**: Verify that users can unfollow other users

**Steps**:
1. Browser 1: Login as User A
2. Browser 2: Login as User B
3. Ensure User B is following User A
4. With User B, navigate to User A's profile
5. Click "Unfollow" button
6. Verify the button changes back to "Follow"
7. Check that User A's posts no longer appear in User B's feed
8. Verify User B is removed from User A's followers list

**Expected Results**:
- ✅ Unfollow action is successful
- ✅ Button changes from "Unfollow" to "Follow"
- ✅ User B no longer sees User A's private/followers-only posts
- ✅ User B is removed from User A's followers list
- ✅ User A's public posts may still be visible

---

## 2. PROFILE TESTS

### Test 2.1: View Own Profile - Complete Information
**Objective**: Verify that user's own profile displays all registration information

**Steps**:
1. Login as User A
2. Navigate to "Profile" or click on your name
3. Verify all information is displayed

**Expected Results**:
- ✅ First Name is displayed
- ✅ Last Name is displayed
- ✅ Email is displayed
- ✅ Date of Birth is displayed
- ✅ Avatar/Profile Picture is displayed (if uploaded)
- ✅ Nickname is displayed (if provided)
- ✅ About Me section is displayed (if provided)
- ✅ Password is NOT displayed (security)
- ✅ Profile privacy status (Public/Private) is shown

---

### Test 2.2: View Own Profile - Posts Display
**Objective**: Verify that user's profile shows all their posts

**Steps**:
1. Login as User A
2. Create 3-5 posts with different privacy settings
3. Navigate to your profile
4. Scroll through the posts section

**Expected Results**:
- ✅ All posts created by User A are displayed
- ✅ Posts are ordered by creation date (newest first)
- ✅ Each post shows content, images (if any), and privacy level
- ✅ Posts with all privacy levels are visible on own profile

---

### Test 2.3: View Own Profile - Followers and Following
**Objective**: Verify that profile displays followers and following lists

**Steps**:
1. Login as User A
2. Ensure User A follows at least 2 users
3. Ensure User A has at least 2 followers
4. Navigate to your profile
5. Check the followers/following section

**Expected Results**:
- ✅ Followers count is accurate
- ✅ Following count is accurate
- ✅ Can click to view list of followers
- ✅ Can click to view list of following
- ✅ Each user in the lists is clickable to view their profile

---

### Test 2.4: Toggle Profile Privacy (Public/Private)
**Objective**: Verify that users can change their profile privacy setting

**Steps**:
1. Login as User A
2. Navigate to profile settings or edit profile
3. Find the privacy toggle (Public/Private)
4. Change from Private to Public
5. Save changes
6. Verify the change is reflected
7. Change back to Private
8. Verify the change is reflected

**Expected Results**:
- ✅ Privacy toggle is available and functional
- ✅ Changes are saved successfully
- ✅ Profile badge shows current privacy status
- ✅ Privacy change affects follow behavior (immediate vs request)

---

### Test 2.5: View Followed Private User Profile
**Objective**: Verify that following a private user grants access to their profile

**Steps**:
1. Browser 1: Login as User A (private profile)
2. Browser 2: Login as User B
3. User B sends follow request to User A
4. User A accepts the request
5. User B navigates to User A's profile

**Expected Results**:
- ✅ User B can view User A's full profile
- ✅ User B can see User A's posts
- ✅ User B can see User A's about section
- ✅ User B can see followers/following counts
- ✅ No "This profile is private" message is shown

---

### Test 2.6: View Non-Followed Private User Profile (Blocked)
**Objective**: Verify that non-followers cannot see private profile content

**Steps**:
1. Browser 1: Login as User A (private profile)
2. Browser 2: Login as User B
3. Ensure User B is NOT following User A
4. User B navigates to User A's profile

**Expected Results**:
- ✅ User B sees basic profile information (name, avatar)
- ✅ User B sees "This profile is private" message
- ✅ User B CANNOT see User A's posts
- ✅ User B CANNOT see User A's about section
- ✅ User B sees a "Follow" button to send request
- ✅ Followers/following counts may be hidden

---

### Test 2.7: View Non-Followed Public User Profile
**Objective**: Verify that anyone can view public profiles

**Steps**:
1. Browser 1: Login as User A (public profile)
2. Browser 2: Login as User B
3. Ensure User B is NOT following User A
4. User B navigates to User A's profile

**Expected Results**:
- ✅ User B can view User A's full profile
- ✅ User B can see User A's public posts
- ✅ User B can see User A's about section
- ✅ User B can see followers/following information
- ✅ User B sees a "Follow" button
- ✅ No privacy restrictions are applied

---

### Test 2.8: View Followed Public User Profile
**Objective**: Verify that followers can view public profiles with full access

**Steps**:
1. Browser 1: Login as User A (public profile)
2. Browser 2: Login as User B
3. User B follows User A
4. User B navigates to User A's profile

**Expected Results**:
- ✅ User B can view User A's complete profile
- ✅ User B can see all of User A's posts (public and followers-only)
- ✅ User B sees "Unfollow" button
- ✅ Full access to profile features

---

## 3. POSTS AND COMMENTS TESTS

### Test 3.1: Create Post After Login
**Objective**: Verify that logged-in users can create posts

**Steps**:
1. Login as User A
2. Navigate to dashboard or home feed
3. Find the "Create Post" section
4. Enter post content: "This is my test post"
5. Click "Post" or "Submit"
6. Verify the post appears in the feed

**Expected Results**:
- ✅ Post creation form is accessible after login
- ✅ Post is created successfully
- ✅ Post appears in user's feed immediately
- ✅ Post appears on user's profile
- ✅ Timestamp is displayed correctly

---

### Test 3.2: Create Post with Image (JPG/PNG)
**Objective**: Verify that posts can include JPG or PNG images

**Steps**:
1. Login as User A
2. Navigate to post creation
3. Enter post content
4. Click "Add Image" or image upload button
5. Select a JPG image from your computer
6. Submit the post
7. Repeat with a PNG image

**Expected Results**:
- ✅ Image upload button/field is available
- ✅ JPG images can be uploaded successfully
- ✅ PNG images can be uploaded successfully
- ✅ Image preview is shown before posting
- ✅ Image is displayed in the post after creation
- ✅ Image is properly sized/formatted

---

### Test 3.3: Create Post with GIF
**Objective**: Verify that posts can include GIF images

**Steps**:
1. Login as User A
2. Navigate to post creation
3. Enter post content
4. Upload a GIF file
5. Submit the post

**Expected Results**:
- ✅ GIF files can be uploaded
- ✅ GIF animation plays in the post
- ✅ GIF is displayed correctly in feed and profile

---

### Test 3.4: Create Comment with Image/GIF
**Objective**: Verify that comments can include images and GIFs

**Steps**:
1. Login as User A
2. Find an existing post
3. Click "Comment" button
4. Enter comment text
5. Upload an image (JPG/PNG/GIF)
6. Submit the comment
7. Verify the comment appears with the image

**Expected Results**:
- ✅ Comment creation includes image upload option
- ✅ JPG/PNG images can be added to comments
- ✅ GIF images can be added to comments
- ✅ Images display correctly in comments
- ✅ Comments with images appear under the post

---

### Test 3.5: Post Privacy Settings
**Objective**: Verify that users can specify post privacy levels

**Steps**:
1. Login as User A
2. Create a new post
3. Find the privacy dropdown/selector
4. Verify available options:
   - Public
   - Private
   - Almost Private (Followers Only)
5. Select "Public" and create post
6. Create another post with "Private" setting
7. Create another post with "Almost Private" setting

**Expected Results**:
- ✅ Privacy selector is available during post creation
- ✅ Three privacy options are available
- ✅ Selected privacy is saved with the post
- ✅ Privacy badge/indicator shows on each post
- ✅ Privacy settings affect post visibility

---

### Test 3.6: Almost Private Post - Specify Allowed Users
**Objective**: Verify that "Almost Private" posts allow selecting specific users

**Steps**:
1. Login as User A
2. Create a new post
3. Select "Almost Private" privacy option
4. Look for user selection interface
5. Select specific users (e.g., User B and User C)
6. Submit the post
7. Verify only selected users can see the post

**Expected Results**:
- ✅ User selection interface appears for "Almost Private"
- ✅ Can search and select specific users
- ✅ Can select multiple users
- ✅ Selected users can view the post
- ✅ Non-selected users cannot view the post
- ✅ Post shows list of allowed users (for post creator)

---

## 4. GROUPS TESTS

### Test 4.1: Create Group and Invite Follower
**Objective**: Verify that users can create groups and invite followers

**Steps**:
1. Browser 1: Login as User A
2. Browser 2: Login as User B
3. Ensure User B is following User A
4. User A creates a new group:
   - Group name: "Test Group"
   - Description: "This is a test group"
5. User A invites User B to the group
6. Verify invitation is sent

**Expected Results**:
- ✅ Group creation form is accessible
- ✅ Can enter group name and description
- ✅ Group is created successfully
- ✅ Can invite followers to the group
- ✅ Invitation is sent to User B
- ✅ User A becomes group owner/admin

---

### Test 4.2: Receive and Accept/Decline Group Invitation
**Objective**: Verify that invited users receive and can respond to invitations

**Steps**:
1. Continue from Test 4.1
2. Browser 2 (User B): Check notifications
3. Verify group invitation is received
4. Click "Accept" on the invitation
5. Verify User B joins the group
6. Repeat with User C, but click "Decline"

**Expected Results**:
- ✅ User B receives group invitation notification
- ✅ Invitation shows group name and inviter
- ✅ Can accept invitation (joins group)
- ✅ Can decline invitation (does not join)
- ✅ Accepted users appear in group members list
- ✅ Declined users do not join the group

---

### Test 4.3: Request to Join Group
**Objective**: Verify that users can request to join groups

**Steps**:
1. Browser 1: Login as User A (group owner)
2. Browser 2: Login as User B
3. User B discovers User A's group
4. User B clicks "Request to Join" button
5. User A receives the join request
6. User A accepts the request
7. Verify User B joins the group

**Expected Results**:
- ✅ "Request to Join" button is available for non-members
- ✅ Join request is sent successfully
- ✅ Group owner receives the request notification
- ✅ Owner can accept or decline the request
- ✅ Accepted users join the group
- ✅ Declined users remain non-members

---

### Test 4.4: Non-Owner Member Can Invite Users
**Objective**: Verify that group members (not just owner) can invite others

**Steps**:
1. Browser 1: Login as User A (group owner)
2. Browser 2: Login as User B (group member)
3. Browser 3: Login as User C (not in group)
4. User B (member, not owner) invites User C to the group
5. Verify invitation is sent
6. User C accepts invitation

**Expected Results**:
- ✅ Group members can access invite functionality
- ✅ Members can invite their followers
- ✅ Invitations from members work same as owner invitations
- ✅ Invited users can join the group
- ✅ Group owner may receive notification of new member

---

### Test 4.5: Create Posts in Group
**Objective**: Verify that group members can create posts within the group

**Steps**:
1. Login as User A (group member)
2. Navigate to the group page
3. Find "Create Post" section within the group
4. Create a post: "This is a group post"
5. Submit the post
6. Verify post appears in group feed

**Expected Results**:
- ✅ Group members can create posts in the group
- ✅ Post creation interface is available in group
- ✅ Posts appear in group feed
- ✅ Posts are visible to all group members
- ✅ Posts show author information

---

### Test 4.6: Comment on Group Posts
**Objective**: Verify that group members can comment on group posts

**Steps**:
1. Browser 1: Login as User A (group member)
2. Browser 2: Login as User B (group member)
3. User A creates a post in the group
4. User B views the post
5. User B adds a comment: "Great post!"
6. Verify comment appears under the post

**Expected Results**:
- ✅ Group members can comment on group posts
- ✅ Comments appear under posts
- ✅ Comments show author information
- ✅ All group members can see comments
- ✅ Comment timestamps are displayed

---

### Test 4.7: Create Event in Group
**Objective**: Verify that group events can be created with required fields

**Steps**:
1. Login as User A (group member/owner)
2. Navigate to the group
3. Click "Create Event" button
4. Fill in event details:
   - Title: "Group Meetup"
   - Description: "Let's meet at the park"
   - Date: [Select future date]
   - Time: [Select time]
5. Verify response options are available:
   - Going
   - Not Going
   - (Optional: Maybe)
6. Submit the event

**Expected Results**:
- ✅ Event creation form is accessible
- ✅ Required fields: Title, Description, Date/Time
- ✅ At least two response options: Going, Not Going
- ✅ Event is created successfully
- ✅ Event appears in group events section
- ✅ All group members can see the event

---

### Test 4.8: View and Respond to Group Event
**Objective**: Verify that group members can view and respond to events

**Steps**:
1. Browser 1: Login as User A (created event in Test 4.7)
2. Browser 2: Login as User B (group member)
3. User B navigates to the group
4. User B views the event "Group Meetup"
5. User B clicks "Going" option
6. User A views the event responses
7. Verify User B's response is recorded

**Expected Results**:
- ✅ All group members can see the event
- ✅ Event shows title, description, date, time
- ✅ Response options are clickable
- ✅ User can select "Going" or "Not Going"
- ✅ Response is saved and displayed
- ✅ Event creator can see who responded
- ✅ Response counts are updated (e.g., "5 Going, 2 Not Going")

---

## 5. ADDITIONAL TESTS

### Test 5.1: Search Users
**Objective**: Verify user search functionality

**Steps**:
1. Login as User A
2. Navigate to search page
3. Search for "John" (or another user's name)
4. Verify search results appear
5. Click on a user to view their profile

**Expected Results**:
- ✅ Search functionality is available
- ✅ Can search by name, email, or nickname
- ✅ Search results display matching users
- ✅ Can click users to view profiles

---

### Test 5.2: Notifications System
**Objective**: Verify that notifications work for various actions

**Steps**:
1. Browser 1: Login as User A
2. Browser 2: Login as User B
3. User B performs actions:
   - Sends follow request to User A
   - Comments on User A's post
   - Invites User A to a group
4. User A checks notifications

**Expected Results**:
- ✅ Notifications appear for follow requests
- ✅ Notifications appear for comments
- ✅ Notifications appear for group invites
- ✅ Notifications show relevant information
- ✅ Can click notifications to view details
- ✅ Can mark notifications as read

---

### Test 5.3: Feed Display
**Objective**: Verify that the feed shows appropriate posts

**Steps**:
1. Login as User A
2. Follow several users with different privacy settings
3. View the main feed/dashboard
4. Verify posts from followed users appear

**Expected Results**:
- ✅ Feed shows posts from followed users
- ✅ Public posts from followed users appear
- ✅ Followers-only posts from followed users appear
- ✅ Private posts do NOT appear (unless specifically shared)
- ✅ Posts are ordered by date (newest first)
- ✅ Own posts appear in feed

---

## 6. SECURITY AND EDGE CASES

### Test 6.1: Unauthenticated Access
**Objective**: Verify that unauthenticated users cannot access protected features

**Steps**:
1. Open browser in incognito/private mode
2. Try to access dashboard directly
3. Try to access profile pages
4. Try to create posts

**Expected Results**:
- ✅ Redirected to login page
- ✅ Cannot access protected routes
- ✅ Cannot perform authenticated actions

---

### Test 6.2: Session Persistence
**Objective**: Verify that sessions persist across page refreshes

**Steps**:
1. Login as User A
2. Navigate to dashboard
3. Refresh the page
4. Verify still logged in

**Expected Results**:
- ✅ Session persists after refresh
- ✅ User remains logged in
- ✅ No need to re-authenticate

---

### Test 6.3: Logout Functionality
**Objective**: Verify that logout properly ends the session

**Steps**:
1. Login as User A
2. Navigate to dashboard
3. Click "Logout"
4. Try to access dashboard again

**Expected Results**:
- ✅ Logout is successful
- ✅ Session is terminated
- ✅ Redirected to login/home page
- ✅ Cannot access protected routes after logout

---

## TEST EXECUTION CHECKLIST

### Before Testing
- [ ] Backend server is running on port 8080
- [ ] Frontend server is running on port 3000
- [ ] Database is initialized with migrations
- [ ] At least 2 browsers are available
- [ ] Test user accounts are created

### During Testing
- [ ] Document any bugs or issues found
- [ ] Take screenshots of failures
- [ ] Note any unexpected behavior
- [ ] Test on different browsers

### After Testing
- [ ] Compile list of bugs/issues
- [ ] Verify all critical features work
- [ ] Document any missing features
- [ ] Create bug reports for developers

---

## EXPECTED OUTCOMES SUMMARY

### ✅ MUST WORK
1. User registration and login
2. Profile viewing (own and others)
3. Follow/unfollow with request system
4. Public vs Private profile behavior
5. Post creation with privacy settings
6. Comment creation
7. Image/GIF uploads
8. Group creation and membership
9. Group invitations and requests
10. Group posts and comments
11. Event creation and responses
12. Notifications for key actions

### ⚠️ SHOULD WORK
1. User search
2. Feed filtering by privacy
3. Session management
4. Real-time updates (WebSocket)
5. Profile editing
6. Almost Private post user selection

### 📝 NICE TO HAVE
1. Emoji support in posts/comments
2. Post editing/deletion
3. Group chat
4. Private messaging
5. Advanced notification settings

---

## REPORTING BUGS

When reporting bugs, include:
1. **Test Case**: Which test was being performed
2. **Steps to Reproduce**: Exact steps taken
3. **Expected Result**: What should have happened
4. **Actual Result**: What actually happened
5. **Screenshots**: Visual evidence of the issue
6. **Browser/Environment**: Browser version, OS, etc.
7. **Severity**: Critical, High, Medium, Low

---

## CONCLUSION

This testing plan covers all the features mentioned in your requirements. Execute each test systematically and document the results. The application should pass all "MUST WORK" tests before being considered production-ready.

Good luck with testing! 🚀
