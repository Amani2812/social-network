# Social Network Testing Checklist

Use this checklist to track your testing progress. Mark items as complete as you test them.

## Setup
- [ ] Backend running on port 8080
- [ ] Frontend running on port 3000
- [ ] Database initialized
- [ ] 2+ browsers available (Chrome, Firefox, etc.)
- [ ] Test users created

---

## 1. FOLLOW/UNFOLLOW TESTS

### Private User Follow Requests
- [ ] Can send follow request to private user
- [ ] Request shows as "Pending" for requester
- [ ] Private user receives follow request notification
- [ ] Private user can accept follow request
- [ ] Private user can decline follow request

### Public User Follow
- [ ] Can follow public user immediately (no request)
- [ ] Button changes to "Unfollow" instantly
- [ ] Follower appears in public user's followers list

### Unfollow
- [ ] Can unfollow a user successfully
- [ ] Button changes from "Unfollow" to "Follow"
- [ ] User removed from followers list
- [ ] Posts no longer appear in feed

---

## 2. PROFILE TESTS

### Own Profile
- [ ] Displays first name
- [ ] Displays last name
- [ ] Displays email
- [ ] Displays date of birth
- [ ] Displays avatar (if uploaded)
- [ ] Displays nickname (if provided)
- [ ] Displays about me (if provided)
- [ ] Does NOT display password
- [ ] Shows all user's posts
- [ ] Shows followers count
- [ ] Shows following count
- [ ] Can view followers list
- [ ] Can view following list
- [ ] Shows privacy status (Public/Private)

### Profile Privacy Toggle
- [ ] Can change from Private to Public
- [ ] Can change from Public to Private
- [ ] Changes are saved
- [ ] Privacy badge updates

### Viewing Other Profiles
- [ ] Can view followed private user's full profile
- [ ] Can view followed private user's posts
- [ ] CANNOT view non-followed private user's posts
- [ ] See "This profile is private" message for non-followed private users
- [ ] Can view non-followed public user's full profile
- [ ] Can view non-followed public user's posts
- [ ] Can view followed public user's full profile

---

## 3. POSTS & COMMENTS TESTS

### Post Creation
- [ ] Can create post after login
- [ ] Post appears in feed immediately
- [ ] Post appears on user profile
- [ ] Can upload JPG image to post
- [ ] Can upload PNG image to post
- [ ] Can upload GIF to post
- [ ] Images display correctly in posts

### Post Privacy
- [ ] Can select "Public" privacy
- [ ] Can select "Private" privacy
- [ ] Can select "Almost Private" (Followers Only) privacy
- [ ] Privacy badge shows on posts
- [ ] Can specify allowed users for "Almost Private" posts
- [ ] Only allowed users can see "Almost Private" posts

### Comments
- [ ] Can comment on posts
- [ ] Can upload JPG image to comment
- [ ] Can upload PNG image to comment
- [ ] Can upload GIF to comment
- [ ] Comments display correctly under posts
- [ ] Comment author and timestamp shown

---

## 4. GROUPS TESTS

### Group Creation & Invitations
- [ ] Can create a group
- [ ] Can enter group name and description
- [ ] Can invite followers to group
- [ ] Invited users receive notification
- [ ] Invited users can accept invitation
- [ ] Invited users can decline invitation
- [ ] Accepted users join the group
- [ ] Declined users don't join

### Group Join Requests
- [ ] Can request to join a group
- [ ] Group owner receives join request
- [ ] Owner can accept join request
- [ ] Owner can decline join request
- [ ] Accepted users join the group

### Group Membership
- [ ] Non-owner members can invite users
- [ ] Group members can create posts in group
- [ ] Group posts visible to all members
- [ ] Group members can comment on group posts
- [ ] Comments visible to all members

### Group Events
- [ ] Can create event in group
- [ ] Event requires title
- [ ] Event requires description
- [ ] Event requires date/time
- [ ] Event has "Going" option
- [ ] Event has "Not Going" option
- [ ] All group members can see event
- [ ] Members can respond to event (Going/Not Going)
- [ ] Event creator can see responses
- [ ] Response counts are displayed

---

## 5. ADDITIONAL FEATURES

### Search
- [ ] Can search for users
- [ ] Search results display correctly
- [ ] Can click user to view profile

### Notifications
- [ ] Receive notification for follow requests
- [ ] Receive notification for comments
- [ ] Receive notification for group invites
- [ ] Receive notification for event responses
- [ ] Can view notifications list
- [ ] Can mark notifications as read

### Feed
- [ ] Feed shows posts from followed users
- [ ] Feed shows own posts
- [ ] Feed respects privacy settings
- [ ] Posts ordered by date (newest first)
- [ ] Public posts from followed users appear
- [ ] Followers-only posts from followed users appear
- [ ] Private posts don't appear (unless shared)

### Authentication
- [ ] Can register new account
- [ ] Can login with credentials
- [ ] Can logout
- [ ] Session persists after page refresh
- [ ] Cannot access protected routes when logged out
- [ ] Redirected to login when accessing protected routes

---

## 6. BROWSER COMPATIBILITY

### Chrome
- [ ] All features work in Chrome
- [ ] UI displays correctly
- [ ] No console errors

### Firefox
- [ ] All features work in Firefox
- [ ] UI displays correctly
- [ ] No console errors

### Edge (Optional)
- [ ] All features work in Edge
- [ ] UI displays correctly
- [ ] No console errors

---

## 7. BUGS FOUND

### Critical Bugs
1. _______________________________________________
2. _______________________________________________
3. _______________________________________________

### High Priority Bugs
1. _______________________________________________
2. _______________________________________________
3. _______________________________________________

### Medium Priority Bugs
1. _______________________________________________
2. _______________________________________________
3. _______________________________________________

### Low Priority Bugs
1. _______________________________________________
2. _______________________________________________
3. _______________________________________________

---

## 8. MISSING FEATURES

List any features from the requirements that are not implemented:

1. _______________________________________________
2. _______________________________________________
3. _______________________________________________
4. _______________________________________________
5. _______________________________________________

---

## TESTING SUMMARY

**Total Tests**: _____ / _____
**Passed**: _____
**Failed**: _____
**Blocked**: _____

**Overall Status**: ⬜ Pass | ⬜ Fail | ⬜ Needs Work

**Tester Name**: _____________________
**Date**: _____________________
**Notes**: 
_________________________________________________
_________________________________________________
_________________________________________________
