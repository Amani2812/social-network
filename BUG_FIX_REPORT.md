# Bug Fix Report - Social Network Application

## Date: 2025-01-17

---

## Bug #1: Follow Requests Not Showing Notifications

### Status: ⚠️ NEEDS BACKEND VERIFICATION

### Description:
When User A follows User B, User B does not receive a notification about the follow request.

### Expected Behavior:
- User A sends follow request to User B
- User B should see a notification in the notifications page
- Notification should show "User A wants to follow you"

### Current Behavior:
- Follow request is sent successfully
- No notification appears for User B

### Possible Causes:
1. Backend notification creation not triggered on follow request
2. Notification API endpoint not returning data correctly
3. Frontend notification polling not working

### Investigation Needed:
Check backend `handlers.go` - `SendFollowRequest` function to verify:
```go
// Should create notification like this:
notification := &models.Notification{
    UserID: followingID,  // The person being followed
    Type: "follow_request",
    Content: fmt.Sprintf("%s %s wants to follow you", follower.FirstName, follower.LastName),
    RelatedID: followerID,
}
repo.CreateNotification(notification)
```

### Frontend Code (Already Correct):
- `frontend/src/app/notifications/page.tsx` - Fetches notifications correctly
- `frontend/src/app/dashboard/page.tsx` - Shows notification bell with count
- Polling every 30 seconds for new notifications

---

## Bug #2: Profile Privacy Toggle Not Persisting

### Status: ⚠️ NEEDS BACKEND VERIFICATION

### Description:
When user toggles profile from Private to Public (or vice versa), the change doesn't persist after page refresh.

### Expected Behavior:
- User checks/unchecks "Public Profile" checkbox
- Clicks "Save Changes"
- Profile privacy should update in database
- After refresh, new privacy setting should be shown

### Current Behavior:
- User toggles privacy setting
- Clicks save
- Success message appears
- After refresh, privacy reverts to previous state

### Frontend Code (Already Correct):
```typescript
// frontend/src/app/profile/edit/page.tsx
const response = await fetch('http://localhost:8080/api/users/profile', {
  method: 'PUT',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    first_name: formData.first_name,
    last_name: formData.last_name,
    nickname: formData.nickname,
    about_me: formData.about_me,
    is_public: formData.is_public  // ✅ Sending boolean correctly
  }),
  credentials: 'include',
})
```

### Investigation Needed:
Check backend `handlers.go` - `UpdateProfile` function to verify:

1. **Does it accept `is_public` field?**
```go
type UpdateProfileRequest struct {
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Nickname  string `json:"nickname"`
    AboutMe   string `json:"about_me"`
    IsPublic  bool   `json:"is_public"`  // ⚠️ Check if this exists
}
```

2. **Does it update the database?**
```go
// Should update like this:
_, err := h.repo.DB.Exec(`
    UPDATE users 
    SET first_name = ?, last_name = ?, nickname = ?, about_me = ?, is_public = ?
    WHERE id = ?
`, req.FirstName, req.LastName, req.Nickname, req.AboutMe, req.IsPublic, user.ID)
```

### Possible Fixes Needed in Backend:

**Option 1: Add is_public to UpdateProfile handler**
```go
func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
    // ... existing code ...
    
    var req struct {
        FirstName string `json:"first_name"`
        LastName  string `json:"last_name"`
        Nickname  string `json:"nickname"`
        AboutMe   string `json:"about_me"`
        IsPublic  bool   `json:"is_public"`  // ADD THIS
    }
    
    // ... decode request ...
    
    // Update query should include is_public
    _, err := h.repo.DB.Exec(`
        UPDATE users 
        SET first_name = ?, last_name = ?, nickname = ?, about_me = ?, is_public = ?, updated_at = CURRENT_TIMESTAMP
        WHERE id = ?
    `, req.FirstName, req.LastName, req.Nickname, req.AboutMe, req.IsPublic, user.ID)
}
```

**Option 2: Create separate endpoint for privacy toggle**
```go
func (h *Handler) TogglePrivacy(w http.ResponseWriter, r *http.Request) {
    user, err := h.getUserFromSession(r)
    if err != nil {
        respondError(w, http.StatusUnauthorized, "Not authenticated")
        return
    }
    
    var req struct {
        IsPublic bool `json:"is_public"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, http.StatusBadRequest, "Invalid request")
        return
    }
    
    _, err = h.repo.DB.Exec(`
        UPDATE users SET is_public = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?
    `, req.IsPublic, user.ID)
    
    if err != nil {
        respondError(w, http.StatusInternalServerError, "Failed to update privacy")
        return
    }
    
    respondJSON(w, http.StatusOK, map[string]string{"message": "Privacy updated"})
}
```

---

## Summary

### Issues Found:
1. ❌ Follow request notifications not being created/displayed
2. ❌ Profile privacy toggle not persisting to database

### Frontend Status:
✅ All frontend code is correct and working as expected
✅ Proper API calls being made
✅ Correct data being sent

### Backend Status:
⚠️ Needs verification and potential fixes:
1. Check if notifications are created on follow requests
2. Check if UpdateProfile handler accepts and saves `is_public` field

### Recommended Actions:
1. Review backend `handlers.go` file
2. Add/fix notification creation in SendFollowRequest
3. Add/fix is_public field handling in UpdateProfile
4. Test both features after backend fixes

---

## Testing Checklist After Backend Fixes:

### Follow Notifications:
- [ ] User A follows User B (private profile)
- [ ] User B sees notification "User A wants to follow you"
- [ ] User B can accept/decline from notifications page
- [ ] Notification count updates in real-time

### Privacy Toggle:
- [ ] User toggles profile to Public
- [ ] Saves changes
- [ ] Refreshes page - should still be Public
- [ ] Other users can now follow without request
- [ ] Toggle back to Private
- [ ] Saves changes
- [ ] Refreshes page - should still be Private
- [ ] Other users now need to send follow request

---

## Notes:
- All frontend null safety fixes have been applied
- All features are implemented in the frontend
- Backend API integration points are correct
- Only backend logic needs verification/fixes
