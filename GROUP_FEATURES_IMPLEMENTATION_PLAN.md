# Group Features Implementation Plan

This document outlines the step-by-step plan to implement all missing group features identified in the testing guide.

---

## Overview of Missing Features

### Critical (Must Have):
1. ✅ Invite users to groups UI
2. ✅ Request to join group UI
3. ✅ Accept/Decline group invitations in notifications
4. ✅ View and respond to join requests (admin)
5. ✅ Fix event notification navigation

### Important (Should Have):
6. ✅ Display group members list
7. ✅ Show event responses (who's going/not going)
8. ✅ Leave group functionality
9. ✅ Comments on group posts

### Nice to Have:
10. ✅ Remove members (admin)
11. ✅ Edit group details (admin)
12. ✅ Delete group (creator)

---

## Implementation Details

### 1. Invite Users to Groups UI

**File**: `frontend/src/app/groups/[id]/page.tsx`

**Changes Needed**:
- Add "Invite Members" button in group header
- Create modal to show list of followers
- Add invite functionality for each follower
- Show invitation status (already invited, already member)

**New State Variables**:
```typescript
const [showInviteModal, setShowInviteModal] = useState(false)
const [followers, setFollowers] = useState<User[]>([])
const [groupMembers, setGroupMembers] = useState<User[]>([])
const [inviting, setInviting] = useState<number | null>(null)
```

**New Functions**:
```typescript
const fetchFollowers = async () => {
  // Fetch current user's followers
}

const fetchGroupMembers = async () => {
  // Fetch group members to check who's already in
}

const handleInviteUser = async (userId: number) => {
  // Call POST /api/groups/invite
}
```

**UI Components**:
- Button: "👥 Invite Members" (visible only to group members)
- Modal with list of followers
- Each follower shows: Avatar, Name, and "Invite" button
- Disable button if user is already a member or invited

---

### 2. Request to Join Group UI

**File**: `frontend/src/app/groups/[id]/page.tsx`

**Changes Needed**:
- Check if current user is a member
- Show "Request to Join" button for non-members
- Handle join request submission
- Show "Request Pending" status after submission

**New State Variables**:
```typescript
const [isMember, setIsMember] = useState(false)
const [hasJoinRequest, setHasJoinRequest] = useState(false)
const [requestingJoin, setRequestingJoin] = useState(false)
```

**New Functions**:
```typescript
const checkMembership = async () => {
  // Check if user is member of the group
}

const handleJoinRequest = async () => {
  // Call POST /api/groups/join/request
}
```

**UI Logic**:
```typescript
{!isMember && !hasJoinRequest && (
  <button onClick={handleJoinRequest}>
    Request to Join
  </button>
)}
{hasJoinRequest && (
  <button disabled>Request Pending</button>
)}
```

---

### 3. Accept/Decline Group Invitations in Notifications

**File**: `frontend/src/app/notifications/page.tsx`

**Changes Needed**:
- Add handler for group_invite notifications
- Add Accept/Decline buttons similar to follow requests
- Call group invite response API
- Refresh notifications after response

**New Functions**:
```typescript
const handleAcceptGroupInvite = async (notification: Notification, e: React.MouseEvent) => {
  e.stopPropagation()
  
  if (!notification.related_id) return
  
  try {
    const response = await fetch('http://localhost:8080/api/groups/invite/respond', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        group_id: notification.related_id,
        accept: true,
      }),
      credentials: 'include',
    })
    
    if (response.ok) {
      await markAsRead(notification.id)
      await fetchNotifications()
      alert('Group invitation accepted! ✓')
    }
  } catch (err) {
    console.error('Failed to accept invitation:', err)
  }
}

const handleDeclineGroupInvite = async (notification: Notification, e: React.MouseEvent) => {
  // Similar to accept but with accept: false
}
```

**UI Update**:
```typescript
{notification.type === 'group_invite' && !notification.is_read && (
  <div className="flex space-x-2 mt-3">
    <button onClick={(e) => handleAcceptGroupInvite(notification, e)}>
      ✓ Accept
    </button>
    <button onClick={(e) => handleDeclineGroupInvite(notification, e)}>
      ✗ Decline
    </button>
  </div>
)}
```

---

### 4. View and Respond to Join Requests (Admin)

**File**: `frontend/src/app/groups/[id]/page.tsx`

**Changes Needed**:
- Add "Join Requests" tab for admins
- Fetch pending join requests
- Display list of users requesting to join
- Add Accept/Decline buttons for each request

**New State Variables**:
```typescript
const [joinRequests, setJoinRequests] = useState<JoinRequest[]>([])
const [isAdmin, setIsAdmin] = useState(false)
const [activeTab, setActiveTab] = useState<'posts' | 'events' | 'requests'>('posts')
```

**New Interface**:
```typescript
interface JoinRequest {
  id: number
  group_id: number
  user_id: number
  status: string
  role: string
  created_at: string
  user?: User
}
```

**New Functions**:
```typescript
const fetchJoinRequests = async () => {
  try {
    const response = await fetch(
      `http://localhost:8080/api/groups/join/requests?group_id=${params.id}`,
      { credentials: 'include' }
    )
    if (response.ok) {
      const data = await response.json()
      setJoinRequests(data || [])
    }
  } catch (err) {
    console.error('Failed to fetch join requests:', err)
  }
}

const handleRespondToJoinRequest = async (userId: number, accept: boolean) => {
  try {
    const response = await fetch('http://localhost:8080/api/groups/join/respond', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        group_id: parseInt(params.id as string),
        user_id: userId,
        accept: accept,
      }),
      credentials: 'include',
    })
    
    if (response.ok) {
      await fetchJoinRequests()
      alert(accept ? 'Request accepted! ✓' : 'Request declined ✗')
    }
  } catch (err) {
    console.error('Failed to respond to request:', err)
  }
}
```

**UI Components**:
```typescript
{isAdmin && (
  <button onClick={() => setActiveTab('requests')}>
    📋 Join Requests {joinRequests.length > 0 && `(${joinRequests.length})`}
  </button>
)}

{activeTab === 'requests' && (
  <div className="p-6">
    {joinRequests.length === 0 ? (
      <div>No pending requests</div>
    ) : (
      joinRequests.map((request) => (
        <div key={request.id}>
          <User info />
          <button onClick={() => handleRespondToJoinRequest(request.user_id, true)}>
            Accept
          </button>
          <button onClick={() => handleRespondToJoinRequest(request.user_id, false)}>
            Decline
          </button>
        </div>
      ))
    )}
  </div>
)}
```

---

### 5. Fix Event Notification Navigation

**File**: `frontend/src/app/notifications/page.tsx`

**Current Code**:
```typescript
} else if (notification.type === 'event_invite' && notification.related_id) {
  router.push(`/events/${notification.related_id}`)
}
```

**Fixed Code**:
```typescript
} else if (notification.type === 'event_invite' && notification.related_id) {
  // Get the event to find its group_id
  try {
    const response = await fetch(
      `http://localhost:8080/api/events/get?id=${notification.related_id}`,
      { credentials: 'include' }
    )
    if (response.ok) {
      const event = await response.json()
      router.push(`/groups/${event.group_id}?tab=events`)
    }
  } catch (err) {
    console.error('Failed to get event:', err)
  }
}
```

**Alternative Simpler Fix**:
Store group_id in notification.related_id instead of event_id when creating event notifications.

**Backend Change** (`backend/pkg/handlers/handlers.go`):
```go
// In CreateEvent handler, change:
h.repo.CreateNotification(member.UserID, "event_invite", content, &event.ID)
// To:
h.repo.CreateNotification(member.UserID, "event_invite", content, &req.GroupID)
```

---

### 6. Display Group Members List

**File**: `frontend/src/app/groups/[id]/page.tsx`

**Changes Needed**:
- Add "Members" tab
- Fetch and display group members
- Show member roles (admin/member)
- Show member count in header

**New State Variables**:
```typescript
const [members, setMembers] = useState<GroupMember[]>([])
const [activeTab, setActiveTab] = useState<'posts' | 'events' | 'members'>('posts')
```

**New Interface**:
```typescript
interface GroupMember {
  id: number
  group_id: number
  user_id: number
  status: string
  role: string
  created_at: string
  user?: User
}
```

**New Functions**:
```typescript
const fetchMembers = async () => {
  try {
    const response = await fetch(
      `http://localhost:8080/api/groups/members?group_id=${params.id}`,
      { credentials: 'include' }
    )
    if (response.ok) {
      const data = await response.json()
      setMembers(data || [])
    }
  } catch (err) {
    console.error('Failed to fetch members:', err)
  }
}
```

**Backend Addition** (`backend/pkg/handlers/handlers.go`):
```go
func (h *Handler) GetGroupMembersWithUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	groupIDStr := r.URL.Query().Get("group_id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	members, err := h.repo.GetGroupMembers(groupID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get members")
		return
	}

	// Enrich with user data
	type MemberWithUser struct {
		*models.GroupMember
		User *models.User `json:"user"`
	}

	var enrichedMembers []MemberWithUser
	for _, member := range members {
		user, _ := h.repo.GetUserByID(member.UserID)
		enrichedMembers = append(enrichedMembers, MemberWithUser{
			GroupMember: member,
			User:        user,
		})
	}

	respondJSON(w, http.StatusOK, enrichedMembers)
}
```

---

### 7. Show Event Responses

**File**: `frontend/src/app/groups/[id]/page.tsx`

**Changes Needed**:
- Fetch event responses for each event
- Display count of "Going" and "Not Going"
- Show list of users and their responses
- Highlight current user's response

**New Interface**:
```typescript
interface EventWithResponses extends Event {
  responses?: EventResponse[]
  going_count?: number
  not_going_count?: number
  user_response?: string
}

interface EventResponse {
  id: number
  event_id: number
  user_id: number
  response: string
  created_at: string
  user?: User
}
```

**Backend Addition** (`backend/pkg/handlers/handlers.go`):
```go
func (h *Handler) GetEventResponses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	eventIDStr := r.URL.Query().Get("event_id")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	rows, err := h.repo.db.Query(
		"SELECT id, event_id, user_id, response, created_at FROM event_responses WHERE event_id = ?",
		eventID,
	)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get responses")
		return
	}
	defer rows.Close()

	type ResponseWithUser struct {
		ID        int       `json:"id"`
		EventID   int       `json:"event_id"`
		UserID    int       `json:"user_id"`
		Response  string    `json:"response"`
		CreatedAt time.Time `json:"created_at"`
		User      *User     `json:"user"`
	}

	var responses []ResponseWithUser
	for rows.Next() {
		var r ResponseWithUser
		if err := rows.Scan(&r.ID, &r.EventID, &r.UserID, &r.Response, &r.CreatedAt); err != nil {
			continue
		}
		r.User, _ = h.repo.GetUserByID(r.UserID)
		responses = append(responses, r)
	}

	respondJSON(w, http.StatusOK, responses)
}
```

**Frontend Update**:
```typescript
const fetchEventResponses = async (eventId: number) => {
  try {
    const response = await fetch(
      `http://localhost:8080/api/events/responses?event_id=${eventId}`,
      { credentials: 'include' }
    )
    if (response.ok) {
      return await response.json()
    }
  } catch (err) {
    console.error('Failed to fetch event responses:', err)
  }
  return []
}

// In event display:
<div className="mt-3">
  <div className="text-sm text-gray-600">
    <span className="text-green-600">✓ {goingCount} Going</span>
    <span className="ml-4 text-red-600">✗ {notGoingCount} Not Going</span>
  </div>
  {showResponses && (
    <div className="mt-2 space-y-1">
      {responses.map(r => (
        <div key={r.id} className="text-sm">
          {r.user?.first_name} {r.user?.last_name} - {r.response}
        </div>
      ))}
    </div>
  )}
</div>
```

---

### 8. Leave Group Functionality

**File**: `frontend/src/app/groups/[id]/page.tsx`

**Changes Needed**:
- Add "Leave Group" button for members (not creator)
- Confirm before leaving
- Remove user from group
- Redirect to groups page

**Backend Addition** (`backend/pkg/handlers/handlers.go`):
```go
func (h *Handler) LeaveGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	var req struct {
		GroupID int `json:"group_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Check if user is the creator
	group, err := h.repo.GetGroup(req.GroupID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Group not found")
		return
	}

	if group.CreatorID == user.ID {
		respondError(w, http.StatusForbidden, "Creator cannot leave group. Delete the group instead.")
		return
	}

	_, err = h.repo.db.Exec(
		"DELETE FROM group_members WHERE group_id = ? AND user_id = ?",
		req.GroupID, user.ID,
	)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to leave group")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Left group successfully"})
}
```

**Frontend**:
```typescript
const handleLeaveGroup = async () => {
  if (!confirm('Are you sure you want to leave this group?')) return
  
  try {
    const response = await fetch('http://localhost:8080/api/groups/leave', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ group_id: parseInt(params.id as string) }),
      credentials: 'include',
    })
    
    if (response.ok) {
      alert('Left group successfully')
      router.push('/groups')
    }
  } catch (err) {
    console.error('Failed to leave group:', err)
  }
}

// In UI (show only if not creator):
{isMember && group.creator_id !== user?.id && (
  <button onClick={handleLeaveGroup} className="bg-red-600 text-white">
    Leave Group
  </button>
)}
```

---

### 9. Comments on Group Posts

**File**: `frontend/src/app/groups/[id]/page.tsx`

**Changes Needed**:
- Add comment section under each post
- Fetch comments for each post
- Add comment input and submit
- Display comments with user info

**Backend Addition** (`backend/pkg/handlers/handlers.go`):
```go
func (h *Handler) CreateGroupPostComment(w http.ResponseWriter, r *http.Request) {
	// Similar to CreateComment but for group posts
}

func (h *Handler) GetGroupPostComments(w http.ResponseWriter, r *http.Request) {
	// Similar to GetComments but for group posts
}
```

**Database Schema Addition**:
```sql
CREATE TABLE IF NOT EXISTS group_post_comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    image_path TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES group_posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

---

## Implementation Priority Order

### Phase 1: Critical Features (Week 1)
1. ✅ Request to Join Group UI
2. ✅ Accept/Decline Group Invitations in Notifications
3. ✅ Fix Event Notification Navigation
4. ✅ View and Respond to Join Requests (Admin)

### Phase 2: Important Features (Week 2)
5. ✅ Invite Users to Groups UI
6. ✅ Display Group Members List
7. ✅ Leave Group Functionality

### Phase 3: Enhanced Features (Week 3)
8. ✅ Show Event Responses
9. ✅ Comments on Group Posts
10. ✅ Remove Members (Admin)
11. ✅ Edit/Delete Group (Admin/Creator)

---

## Testing After Implementation

After implementing each feature, test using the scenarios in `COMPREHENSIVE_GROUP_TESTING_GUIDE.md`.

### Quick Test Checklist:
- [ ] Can invite users to groups
- [ ] Can request to join groups
- [ ] Can accept/decline group invites from notifications
- [ ] Admins can see and respond to join requests
- [ ] Event notifications navigate correctly
- [ ] Can see list of group members
- [ ] Can see event responses
- [ ] Can leave groups
- [ ] Can comment on group posts

---

## Files to Modify

### Frontend:
1. `frontend/src/app/groups/[id]/page.tsx` - Main group page (most changes)
2. `frontend/src/app/notifications/page.tsx` - Group invite handling
3. `frontend/src/app/groups/page.tsx` - Minor updates for join status

### Backend:
1. `backend/pkg/handlers/handlers.go` - New endpoints
2. `backend/server.go` - Register new routes
3. `backend/pkg/models/models.go` - New methods if needed
4. `backend/pkg/db/database.go` - Schema updates for comments

### New Files:
1. `frontend/src/components/InviteMembersModal.tsx` - Reusable invite modal
2. `frontend/src/components/JoinRequestsList.tsx` - Join requests component
3. `frontend/src/components/GroupMembersList.tsx` - Members list component
4. `frontend/src/components/EventResponses.tsx` - Event responses component

---

## Conclusion

This implementation plan provides a clear roadmap to complete all missing group features. By following this plan in phases, you can systematically add functionality while maintaining code quality and testability.

Each feature is designed to integrate seamlessly with the existing codebase and follows the established patterns in the application.
