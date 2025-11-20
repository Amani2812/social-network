# Implementation TODO List

## Feature 1: Message Notifications for Offline Users
- [x] Modify `backend/pkg/websocket/websocket.go` - Add notification creation when user is offline
- [x] Update `frontend/src/app/notifications/page.tsx` - Add message notification type handling

## Feature 2: Group Discovery (All Groups Visible)
- [x] Add `GetAllGroups` method in `backend/pkg/models/models.go`
- [x] Add `GetAllGroups` handler in `backend/pkg/handlers/handlers.go`
- [x] Add route in `backend/server.go` for `/api/groups/all`
- [x] Update `frontend/src/app/groups/page.tsx` - Add "Discover Groups" section

## Bug Fixes
- [x] Fixed JSON parsing error in messages page WebSocket handler

## Testing
- [ ] Test message notifications when recipient is offline
- [ ] Test group discovery and visibility
- [ ] Verify notification navigation works correctly

## Summary of Changes

### Backend Changes:
1. **websocket.go**: Modified `sendToUser` function to create a notification when the recipient is offline
2. **models.go**: Added `GetAllGroups()` method to fetch all groups from the database
3. **handlers.go**: Added `GetAllGroups()` handler to expose the endpoint
4. **server.go**: Added route `/api/groups/all` for fetching all groups

### Frontend Changes:
1. **notifications/page.tsx**: 
   - Added 'message' case in `getNotificationIcon()` with 💬 icon
   - Added navigation to messages page when clicking message notifications
2. **groups/page.tsx**:
   - Added `allGroups` state to store all available groups
   - Added `fetchAllGroups()` function to fetch all groups
   - Added "Discover Groups" section showing all groups with member badges
   - Groups show "Member" badge if user is already a member
