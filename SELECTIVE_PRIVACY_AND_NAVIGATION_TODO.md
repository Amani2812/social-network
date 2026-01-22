# Implementation Plan: Selective Post Privacy + Global Navigation

## Feature 1: Selective Post Privacy

### Backend Changes

#### 1. Database ✅
- [x] Added `post_allowed_users` table
- [x] Added indexes for performance

#### 2. Models (backend/pkg/models/models.go)
- [ ] Update `CreatePost` to accept `allowedUserIDs []int` parameter
- [ ] Add method `AddAllowedUsers(postID int, userIDs []int) error`
- [ ] Update `GetFeed` to check `post_allowed_users` table for custom privacy
- [ ] Update `GetUserPosts` to check `post_allowed_users` table

#### 3. Handlers (backend/pkg/handlers/handlers.go)
- [ ] Update `CreatePost` handler to accept `allowed_users` array in request
- [ ] Add `GetFollowersForPost` handler to return user's followers for selection

#### 4. Routes (backend/server.go)
- [ ] Add route: `GET /api/users/followers` - Get current user's followers

### Frontend Changes

#### 1. Dashboard (frontend/src/app/dashboard/page.tsx)
- [ ] Update create post form to include privacy selector
- [ ] Add "Custom" privacy option
- [ ] When "Custom" selected, show follower selector (checkboxes/multi-select)
- [ ] Fetch user's followers when custom privacy selected
- [ ] Send `allowed_users` array with post creation

#### 2. Types
- [ ] Update Post interface to include `allowed_users?: number[]`

---

## Feature 2: Global Navigation (Dashboard + Notifications on Every Page)

### Approach: Create a Shared Navigation Component

#### 1. Create Navigation Component (frontend/src/components/Navigation.tsx)
- [ ] Create new file with navigation bar
- [ ] Include links to:
  - Dashboard
  - Messages
  - Groups
  - Profile
  - Notifications (with unread count badge)
- [ ] Add logout button
- [ ] Make it responsive (mobile-friendly)
- [ ] Add WebSocket connection for real-time notification updates

#### 2. Update Layout (frontend/src/app/layout.tsx)
- [ ] Import Navigation component
- [ ] Add Navigation to layout (appears on all pages)
- [ ] Exclude from login/register pages

#### 3. Style Navigation
- [ ] Add CSS for fixed/sticky navigation bar
- [ ] Add notification badge styling
- [ ] Add active link highlighting
- [ ] Ensure it doesn't overlap with page content

---

## Implementation Order

### Phase 1: Selective Privacy (Backend)
1. Update models for allowed users
2. Update handlers to accept allowed users
3. Update feed query to check allowed users
4. Test with API calls

### Phase 2: Selective Privacy (Frontend)
1. Update dashboard create post UI
2. Add follower selector
3. Test post creation with custom privacy

### Phase 3: Global Navigation
1. Create Navigation component
2. Add to layout
3. Test on all pages
4. Add real-time notification updates

---

## Testing Checklist

### Selective Privacy
- [ ] Create post with "Public" privacy
- [ ] Create post with "Private" privacy
- [ ] Create post with "Followers Only" privacy
- [ ] Create post with "Custom" privacy (select specific users)
- [ ] Verify only selected users can see custom privacy posts
- [ ] Verify feed shows correct posts based on privacy

### Global Navigation
- [ ] Navigation appears on all pages
- [ ] Navigation does NOT appear on login/register
- [ ] Notification count updates in real-time
- [ ] Links work correctly
- [ ] Logout works
- [ ] Mobile responsive
- [ ] Active page is highlighted

---

## Estimated Time
- Selective Privacy: 2-3 hours
- Global Navigation: 1-2 hours
- Testing: 1 hour
- **Total: 4-6 hours**

---

## Notes
- Database migration will run automatically on server restart
- Need to restart backend after model/handler changes
- Frontend hot-reloads automatically
- Test with multiple users to verify privacy settings
