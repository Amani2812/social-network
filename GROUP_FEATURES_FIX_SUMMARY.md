# Group Features Fix Summary

## Issues Resolved

### 1. Join Request Display Issue ✅
**Problem**: Join requests showed "User ID: X" instead of actual user names.

**Solution**: 
- Updated backend `GetGroupJoinRequests` to JOIN with users table
- Modified `GroupMember` struct to include `User *User` field
- Frontend now displays full user information with avatar, name, and nickname

### 2. Event Response UI Not Updating ✅
**Problem**: After clicking "Going" or "Not Going", the UI didn't reflect the user's choice.

**Solution**:
- Added `user_response` field to `GroupEvent` struct
- Created `GetGroupEventsWithUserResponse` method that LEFT JOINs with event_responses
- Updated frontend to show and disable the selected response button
- Added visual feedback with "(You)" indicator on selected option

### 3. Duplicate Join Request Error ✅
**Problem**: Users got "Failed to send join request" when trying to join a group they already requested.

**Solution**: The error handling was already in place in the backend. The error message now properly displays to users, informing them they already have a pending request or are already a member.

## Technical Implementation

### Backend Changes

#### `backend/pkg/models/models.go`

**GroupMember Struct Update:**
```go
type GroupMember struct {
    ID        int       `json:"id"`
    GroupID   int       `json:"group_id"`
    UserID    int       `json:"user_id"`
    Status    string    `json:"status"`
    Role      string    `json:"role"`
    CreatedAt time.Time `json:"created_at"`
    User      *User     `json:"user,omitempty"`  // NEW
}
```

**GroupEvent Struct Update:**
```go
type GroupEvent struct {
    ID           int       `json:"id"`
    GroupID      int       `json:"group_id"`
    CreatorID    int       `json:"creator_id"`
    Title        string    `json:"title"`
    Description  *string   `json:"description,omitempty"`
    EventTime    time.Time `json:"event_time"`
    CreatedAt    time.Time `json:"created_at"`
    UserResponse *string   `json:"user_response,omitempty"`  // NEW
}
```

**New Method - GetGroupJoinRequests:**
```go
func (r *Repository) GetGroupJoinRequests(groupID int) ([]*GroupMember, error) {
    rows, err := r.db.Query(
        `SELECT gm.id, gm.group_id, gm.user_id, gm.status, gm.role, gm.created_at,
        u.id, u.email, u.first_name, u.last_name, u.avatar_path, u.nickname
        FROM group_members gm
        JOIN users u ON gm.user_id = u.id
        WHERE gm.group_id = ? AND gm.status = 'pending'`,
        groupID,
    )
    // ... implementation
}
```

**New Method - GetGroupEventsWithUserResponse:**
```go
func (r *Repository) GetGroupEventsWithUserResponse(groupID, userID int) ([]*GroupEvent, error) {
    rows, err := r.db.Query(
        `SELECT ge.id, ge.group_id, ge.creator_id, ge.title, ge.description, ge.event_time, ge.created_at,
        er.response
        FROM group_events ge
        LEFT JOIN event_responses er ON ge.id = er.event_id AND er.user_id = ?
        WHERE ge.group_id = ?
        ORDER BY ge.event_time ASC`,
        userID, groupID,
    )
    // ... implementation
}
```

#### `backend/pkg/handlers/handlers.go`

**Updated GetGroupEvents Handler:**
```go
func (h *Handler) GetGroupEvents(w http.ResponseWriter, r *http.Request) {
    user, err := h.getUserFromSession(r)
    if err != nil {
        respondError(w, http.StatusUnauthorized, "Not authenticated")
        return
    }

    groupIDStr := r.URL.Query().Get("group_id")
    groupID, err := strconv.Atoi(groupIDStr)
    if err != nil {
        respondError(w, http.StatusBadRequest, "Invalid group ID")
        return
    }

    events, err := h.repo.GetGroupEventsWithUserResponse(groupID, user.ID)
    if err != nil {
        respondError(w, http.StatusInternalServerError, "Failed to get events")
        return
    }

    respondJSON(w, http.StatusOK, events)
}
```

### Frontend Changes

#### `frontend/src/app/groups/[id]/page.tsx`

**Updated Event Interface:**
```typescript
interface Event {
  id: number
  group_id: number
  creator_id: number
  title: string
  description?: string
  event_time: string
  created_at: string
  user_response?: string  // NEW
}
```

**New JoinRequest Interface:**
```typescript
interface JoinRequest {
  id: number
  group_id: number
  user_id: number
  status: string
  role: string
  created_at: string
  user?: User  // NEW
}
```

**Updated Join Request Display:**
```tsx
{joinRequests.map((request) => (
  <div key={request.user_id} className="border border-gray-200 rounded-lg p-4 flex items-center justify-between">
    <div className="flex items-center">
      <div className="w-12 h-12 bg-gray-300 rounded-full flex items-center justify-center mr-3">
        <span className="text-lg text-gray-600">
          {request.user?.first_name?.[0]}{request.user?.last_name?.[0]}
        </span>
      </div>
      <div>
        <p className="font-semibold text-gray-900">
          {request.user?.first_name} {request.user?.last_name}
        </p>
        {request.user?.nickname && (
          <p className="text-sm text-gray-600">@{request.user.nickname}</p>
        )}
        <p className="text-sm text-gray-500">
          Requested {new Date(request.created_at).toLocaleString()}
        </p>
      </div>
    </div>
    {/* Accept/Decline buttons */}
  </div>
))}
```

**Updated Event Response Buttons:**
```tsx
<div className="flex space-x-2">
  <button
    onClick={() => handleEventResponse(event.id, 'going')}
    disabled={event.user_response === 'going'}
    className={`px-4 py-1 rounded-md text-sm font-medium ${
      event.user_response === 'going'
        ? 'bg-green-700 text-white cursor-default'
        : 'bg-green-600 hover:bg-green-700 text-white'
    }`}
  >
    {event.user_response === 'going' ? '✓ Going (You)' : '✓ Going'}
  </button>
  <button
    onClick={() => handleEventResponse(event.id, 'not_going')}
    disabled={event.user_response === 'not_going'}
    className={`px-4 py-1 rounded-md text-sm font-medium ${
      event.user_response === 'not_going'
        ? 'bg-red-700 text-white cursor-default'
        : 'bg-red-600 hover:bg-red-700 text-white'
    }`}
  >
    {event.user_response === 'not_going' ? '✗ Not Going (You)' : '✗ Not Going'}
  </button>
</div>
```

## Testing Instructions

### 1. Test Join Request Display
1. Create a new user account
2. Request to join a group
3. Log in as the group admin
4. Navigate to the group's "Manage Members" tab
5. **Expected**: See the requester's full name, avatar initials, and nickname (if set)

### 2. Test Event Response UI
1. As a group member, navigate to the Events tab
2. Click "Going" on an event
3. **Expected**: 
   - Alert shows "You responded: Going ✓"
   - "Going" button changes to darker green with "(You)" label
   - "Going" button becomes disabled
   - "Not Going" button remains clickable
4. Click "Not Going" on another event
5. **Expected**: Similar behavior but for "Not Going" button

### 3. Test Duplicate Join Request Prevention
1. Request to join a group
2. Try to request again
3. **Expected**: Error message "user already has a relationship with this group"

### 4. Test Response Persistence
1. Respond to an event
2. Refresh the page
3. **Expected**: Your response is still shown and the button is still disabled

## Benefits

1. **Better UX**: Users can see who is requesting to join their group
2. **Clear Feedback**: Event responses are immediately visible and persistent
3. **Prevents Confusion**: Disabled buttons prevent accidental duplicate responses
4. **Professional Look**: Avatar initials and nicknames make the interface more personal

## Files Modified

- `backend/pkg/models/models.go` - Added user details to structs and new query methods
- `backend/pkg/handlers/handlers.go` - Updated event handler to include user responses
- `frontend/src/app/groups/[id]/page.tsx` - Enhanced UI with user details and response states
- `TODO.md` - Tracked implementation progress

## Database Schema (No Changes Required)

The existing schema already supports these features:
- `users` table has all necessary user information
- `group_members` table tracks membership and requests
- `event_responses` table stores user responses to events

The fixes only required updating the queries to JOIN these tables appropriately.
