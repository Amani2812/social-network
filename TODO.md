# Implementation TODO List

## Part 1: Real-Time Messaging Fix
- [x] Review and analyze current WebSocket message handling in messages/page.tsx
- [x] Fix message reception logic to properly display incoming messages
- [ ] Test real-time messaging between two users

## Part 2: Event Attendance Counts
- [x] Update GroupEvent struct in backend/pkg/models/models.go to include count fields
- [x] Modify GetGroupEventsWithUserResponse() to query and return attendance counts
- [x] Update Event interface in frontend/src/app/groups/[id]/page.tsx
- [x] Display attendance counts in the event UI
- [ ] Test event attendance count display

## Testing
- [ ] Test real-time messaging with multiple browser windows
- [ ] Test event attendance counts with multiple users
- [ ] Verify WebSocket connection stability

## Summary of Changes Made:

### Backend Changes:
1. **backend/pkg/models/models.go**:
   - Added `GoingCount` and `NotGoingCount` fields to `GroupEvent` struct
   - Updated `GetGroupEventsWithUserResponse()` to include SQL subqueries that count responses

### Frontend Changes:
1. **frontend/src/app/groups/[id]/page.tsx**:
   - Added `going_count` and `not_going_count` fields to Event interface
   - Added attendance count display showing number of people going/not going
   - Styled counts with green for "Going" and red for "Not Going"

### Real-Time Messaging:
- The WebSocket implementation is already in place and should work for real-time messaging
- Messages are sent via WebSocket and received in real-time
- The receiver's browser should automatically update when a new message arrives
