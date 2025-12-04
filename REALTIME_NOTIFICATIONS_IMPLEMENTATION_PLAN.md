# Real-Time Notifications Implementation Plan

## Current State Analysis

### ✅ What's Already Working:
1. **Database notifications** are created for:
   - Follow requests (both public and private profiles)
   - Group invites
   - Group join requests (to admins)
   - Event creation (to all group members)
   - Private messages

2. **WebSocket infrastructure** exists:
   - Hub with client management
   - `SendNotificationToUser()` method available
   - Real-time messaging working
   - Follow status updates working

3. **Frontend notification system**:
   - Notifications page displays all notifications
   - Unread count badge
   - Click handlers for different notification types

### ❌ What's Missing:

1. **Real-time notification delivery via WebSocket**:
   - Group join requests don't trigger real-time notifications
   - Event creation doesn't trigger real-time notifications
   - Notifications are only visible after page refresh

2. **Frontend WebSocket notification handling**:
   - Dashboard receives notifications but limited handling
   - No real-time notification updates on notifications page
   - No toast/popup for new notifications

## Implementation Plan

### Phase 1: Backend - Real-Time Notification Delivery

#### 1.1 Update `RequestToJoinGroup` Handler
**File**: `backend/pkg/handlers/handlers.go`

**Current code** (lines ~1050-1075):
```go
func (h *Handler) RequestToJoinGroup(w http.ResponseWriter, r *http.Request) {
    // ... existing code ...
    
    // Create notification for group admins
    group, _ := h.repo.GetGroup(req.GroupID)
    if group != nil {
        members, _ := h.repo.GetGroupMembers(req.GroupID)
        for _, member := range members {
            if member.Role == "admin" {
                content := fmt.Sprintf("%s %s wants to join %s", user.FirstName, user.LastName, group.Title)
                h.repo.CreateNotification(member.UserID, "group_join_request", content, &req.GroupID)
            }
        }
    }
}
```

**Need to add**: Real-time WebSocket notification after creating database notification
```go
// After h.repo.CreateNotification()
h.hub.SendNotificationToUser(member.UserID, "group_join_request", content, user.ID)
```

#### 1.2 Update `CreateEvent` Handler
**File**: `backend/pkg/handlers/handlers.go`

**Current code** (lines ~1200-1230):
```go
func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
    // ... existing code ...
    
    // Notify all group members
    members, _ := h.repo.GetGroupMembers(req.GroupID)
    group, _ := h.repo.GetGroup(req.GroupID)
    if group != nil {
        for _, member := range members {
            if member.UserID != user.ID {
                content := fmt.Sprintf("New event in %s: %s", group.Title, req.Title)
                h.repo.CreateNotification(member.UserID, "event_invite", content, &event.ID)
            }
        }
    }
}
```

**Need to add**: Real-time WebSocket notification
```go
// After h.repo.CreateNotification()
h.hub.SendNotificationToUser(member.UserID, "event_invite", content, user.ID)
```

#### 1.3 Update Hub Interface
**File**: `backend/pkg/handlers/handlers.go`

**Current interface** (line ~22):
```go
hub interface {
    BroadcastNewPost(postID int, userID int)
    SendFollowStatusUpdate(userID int)
}
```

**Need to add**:
```go
hub interface {
    BroadcastNewPost(postID int, userID int)
    SendFollowStatusUpdate(userID int)
    SendNotificationToUser(userID int, notifType, content string, senderID int)
}
```

### Phase 2: Frontend - Real-Time Notification Reception

#### 2.1 Update Dashboard WebSocket Handler
**File**: `frontend/src/app/dashboard/page.tsx`

**Current code** (lines ~80-95):
```typescript
websocket.onmessage = (event) => {
    // Handle real-time notifications
    if (data.type === 'notification') {
        // Show toast notification
        const notifId = Date.now()
        // Auto-remove notification after 5 seconds
        setTimeout(() => {
            // ...
        }, 5000)
    }
}
```

**Enhancement needed**:
- Refresh unread count when notification received
- Better toast styling
- Sound/visual indicator

#### 2.2 Add WebSocket to Notifications Page
**File**: `frontend/src/app/notifications/page.tsx`

**Need to add**:
- WebSocket connection
- Real-time notification updates
- Auto-refresh notification list when new notification arrives

#### 2.3 Add WebSocket to Groups Page
**File**: `frontend/src/app/groups/[id]/page.tsx`

**Need to add**:
- WebSocket connection for real-time updates
- Listen for event creation notifications
- Auto-refresh events list

### Phase 3: Testing Requirements

#### 3.1 Multi-Browser Testing Setup
**Tools needed**:
- Chrome (User A)
- Firefox/Edge (User B)
- Both browsers in incognito/private mode

#### 3.2 Test Scenarios

**Scenario 1: Group Join Request Notification**
1. User A creates a group (becomes admin)
2. User B requests to join the group
3. **Expected**: User A receives real-time notification
4. **Verify**: 
   - Notification appears in User A's dashboard
   - Unread count updates
   - Notification visible in notifications page

**Scenario 2: Event Creation Notification**
1. User A and User B are both members of a group
2. User A creates an event
3. **Expected**: User B receives real-time notification
4. **Verify**:
   - Notification appears in User B's dashboard
   - Event visible in group events tab
   - Notification clickable to navigate to event

## Implementation Steps

### Step 1: Backend Updates
- [ ] Update Hub interface in handlers.go
- [ ] Add real-time notification to RequestToJoinGroup
- [ ] Add real-time notification to CreateEvent
- [ ] Test backend with curl/Postman

### Step 2: Frontend Updates
- [ ] Enhance dashboard notification handling
- [ ] Add WebSocket to notifications page
- [ ] Add WebSocket to groups page
- [ ] Test real-time updates

### Step 3: Integration Testing
- [ ] Test group join request flow
- [ ] Test event creation flow
- [ ] Test with multiple browsers
- [ ] Verify notification persistence

### Step 4: Documentation
- [ ] Create testing guide
- [ ] Document multi-browser setup
- [ ] Add troubleshooting section

## Files to Modify

### Backend:
1. `backend/pkg/handlers/handlers.go` - Add real-time notifications
2. `backend/pkg/websocket/websocket.go` - Already has SendNotificationToUser (no changes needed)

### Frontend:
1. `frontend/src/app/dashboard/page.tsx` - Enhance notification handling
2. `frontend/src/app/notifications/page.tsx` - Add WebSocket
3. `frontend/src/app/groups/[id]/page.tsx` - Add WebSocket for events

## Success Criteria

✅ User A receives real-time notification when User B requests to join their group
✅ User B receives real-time notification when User A creates an event in shared group
✅ Notifications persist in database
✅ Notifications visible in notifications page
✅ Unread count updates in real-time
✅ Multi-browser testing successful
✅ No page refresh required to see notifications

## Notes

- WebSocket infrastructure already exists and works for messaging
- SendNotificationToUser method already implemented in Hub
- Main work is connecting existing pieces together
- Frontend needs WebSocket connections on more pages
- Backend needs to call SendNotificationToUser after CreateNotification
