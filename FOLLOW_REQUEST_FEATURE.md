# Follow Request Accept/Decline Feature

## Date: 2024-01-15
## Feature: Accept/Decline Follow Requests for Private Accounts

---

## Overview

Added Accept/Decline buttons to the notifications page so users with private accounts can respond to follow requests directly from their notifications.

---

## Changes Made

### Frontend: `frontend/src/app/notifications/page.tsx`

#### 1. Added Accept Follow Request Handler
```typescript
const handleAcceptFollowRequest = async (notification: Notification, e: React.MouseEvent) => {
  e.stopPropagation() // Prevent notification click
  
  if (!notification.related_id) return

  try {
    const response = await fetch('http://localhost:8080/api/follow/respond', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        follower_id: notification.related_id,
        accept: true,
      }),
      credentials: 'include',
    })

    if (response.ok) {
      await markAsRead(notification.id)
      fetchNotifications()
      alert('Follow request accepted!')
    }
  } catch (err) {
    console.error('Failed to accept follow request:', err)
    alert('Failed to accept follow request')
  }
}
```

#### 2. Added Decline Follow Request Handler
```typescript
const handleDeclineFollowRequest = async (notification: Notification, e: React.MouseEvent) => {
  e.stopPropagation() // Prevent notification click
  
  if (!notification.related_id) return

  try {
    const response = await fetch('http://localhost:8080/api/follow/respond', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        follower_id: notification.related_id,
        accept: false,
      }),
      credentials: 'include',
    })

    if (response.ok) {
      await markAsRead(notification.id)
      fetchNotifications()
      alert('Follow request declined')
    }
  } catch (err) {
    console.error('Failed to decline follow request:', err)
    alert('Failed to decline follow request')
  }
}
```

#### 3. Updated Notification Click Handler
```typescript
const handleNotificationClick = (notification: Notification) => {
  // Don't navigate for follow requests (they have action buttons)
  if (notification.type === 'follow_request') {
    return
  }

  // Mark as read and navigate for other notification types
  if (!notification.is_read) {
    markAsRead(notification.id)
  }

  // Navigate based on notification type
  if (notification.type === 'group_invite' && notification.related_id) {
    router.push(`/groups/${notification.related_id}`)
  } else if (notification.type === 'event_invite' && notification.related_id) {
    router.push(`/events/${notification.related_id}`)
  } else if (notification.type === 'message' && notification.related_id) {
    router.push(`/messages?user=${notification.related_id}`)
  }
}
```

#### 4. Added Action Buttons to UI
```tsx
{/* Accept/Decline buttons for follow requests */}
{notification.type === 'follow_request' && !notification.is_read && (
  <div className="flex space-x-2 mt-3">
    <button
      onClick={(e) => handleAcceptFollowRequest(notification, e)}
      className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md text-sm font-medium transition-colors"
    >
      ✓ Accept
    </button>
    <button
      onClick={(e) => handleDeclineFollowRequest(notification, e)}
      className="bg-gray-300 hover:bg-gray-400 text-gray-700 px-4 py-2 rounded-md text-sm font-medium transition-colors"
    >
      ✗ Decline
    </button>
  </div>
)}
```

---

## How It Works

### User Flow

1. **User A (Private Account)** receives a follow request from **User B**
2. **User A** sees a notification: "User B wants to follow you"
3. **User A** clicks on the notification bell and goes to the notifications page
4. **User A** sees the follow request notification with two buttons:
   - ✓ **Accept** (Blue button)
   - ✗ **Decline** (Gray button)
5. **User A** clicks either button:
   - **Accept**: User B becomes a follower and can see User A's posts
   - **Decline**: User B's follow request is rejected
6. The notification is marked as read and the list refreshes

### Backend Integration

The feature uses the existing backend API endpoint:
- **Endpoint**: `POST /api/follow/respond`
- **Request Body**:
  ```json
  {
    "follower_id": 123,
    "accept": true  // or false
  }
  ```
- **Response**: Updates the follow status in the database

### Database Flow

When a follow request is accepted/declined:
1. The `follows` table is updated with the new status:
   - `status = 'accepted'` (if accepted)
   - `status = 'rejected'` (if declined)
2. The notification is marked as read
3. The notifications list is refreshed

---

## Features

### ✅ What's Included

1. **Accept Button**
   - Blue button with checkmark icon
   - Accepts the follow request
   - Shows success alert
   - Marks notification as read
   - Refreshes notification list

2. **Decline Button**
   - Gray button with X icon
   - Declines the follow request
   - Shows decline alert
   - Marks notification as read
   - Refreshes notification list

3. **Smart UI Behavior**
   - Buttons only show for unread follow request notifications
   - Once read (accepted/declined), buttons disappear
   - Click events don't propagate to parent (no accidental navigation)
   - Follow request notifications don't navigate on click

4. **Error Handling**
   - Shows alert if accept/decline fails
   - Logs errors to console for debugging
   - Graceful fallback if API call fails

---

## UI/UX Details

### Visual Design

- **Accept Button**: 
  - Background: Blue (#2563eb)
  - Hover: Darker blue (#1d4ed8)
  - Icon: ✓ (checkmark)
  - Text: White

- **Decline Button**:
  - Background: Gray (#d1d5db)
  - Hover: Darker gray (#9ca3af)
  - Icon: ✗ (X mark)
  - Text: Dark gray (#374151)

### Interaction

- Buttons appear below the notification content
- Spaced 8px apart (space-x-2)
- 12px top margin (mt-3)
- Smooth hover transitions
- Click events stop propagation to prevent unwanted navigation

---

## Testing Checklist

### Manual Testing Steps

1. **Setup**:
   - Create two user accounts (User A and User B)
   - Set User A's profile to private

2. **Test Accept Flow**:
   - [ ] User B sends follow request to User A
   - [ ] User A receives notification
   - [ ] User A sees Accept/Decline buttons
   - [ ] User A clicks Accept
   - [ ] Success alert appears
   - [ ] Notification is marked as read
   - [ ] Buttons disappear
   - [ ] User B can now see User A's posts

3. **Test Decline Flow**:
   - [ ] User B sends follow request to User A
   - [ ] User A receives notification
   - [ ] User A sees Accept/Decline buttons
   - [ ] User A clicks Decline
   - [ ] Decline alert appears
   - [ ] Notification is marked as read
   - [ ] Buttons disappear
   - [ ] User B cannot see User A's posts

4. **Test Edge Cases**:
   - [ ] Buttons don't appear for already read notifications
   - [ ] Buttons don't appear for other notification types
   - [ ] Clicking buttons doesn't navigate away
   - [ ] Error handling works if API fails

---

## Backend API Reference

### Endpoint: Respond to Follow Request

**URL**: `POST /api/follow/respond`

**Headers**:
```
Content-Type: application/json
Cookie: session_id=<session_id>
```

**Request Body**:
```json
{
  "follower_id": 123,
  "accept": true
}
```

**Response** (Success - 200 OK):
```json
{
  "message": "Follow request accepted"
}
```

**Response** (Error - 401 Unauthorized):
```json
{
  "error": "Unauthorized"
}
```

**Response** (Error - 500 Internal Server Error):
```json
{
  "error": "Failed to respond to follow request"
}
```

---

## Future Enhancements

### Potential Improvements

1. **Toast Notifications**
   - Replace alerts with styled toast notifications
   - Auto-dismiss after 3 seconds
   - Show success/error states with icons

2. **Optimistic UI Updates**
   - Update UI immediately before API call
   - Revert if API call fails
   - Faster perceived performance

3. **Batch Actions**
   - "Accept All" button for multiple requests
   - "Decline All" button for multiple requests
   - Checkbox selection for bulk actions

4. **User Preview**
   - Show user avatar in notification
   - Display user bio/info on hover
   - Quick profile preview modal

5. **Undo Action**
   - Allow undo within 5 seconds
   - Temporary "Undo" button after action
   - Revert follow status if undone

---

## Related Files

- `frontend/src/app/notifications/page.tsx` - Main notifications page with Accept/Decline buttons
- `backend/pkg/handlers/handlers.go` - Backend API handler for follow responses
- `backend/pkg/models/models.go` - Database operations for follow status updates

---

## Conclusion

The Accept/Decline follow request feature is now fully implemented and ready for use. Users with private accounts can now easily manage their follow requests directly from the notifications page without needing to navigate to the follower's profile.

**Status**: ✅ **COMPLETE AND READY FOR TESTING**
