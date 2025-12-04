# WebSocket Timeout Fix - TODO

## Problem
WebSocket connections stop working after ~10 minutes without page refresh due to:
- Backend read deadline timeout (120 seconds)
- Client not responding to server pings
- No proper keepalive mechanism

## Solution Steps

### Backend Fixes
- [x] Update backend/pkg/websocket/websocket.go
  - [x] Increase read/write deadlines (60s pongWait, 54s pingPeriod)
  - [x] Improve ping/pong mechanism
  - [x] Add better error handling
  - [x] Add configurable timeout constants

### Frontend Fixes
- [x] Update frontend/src/app/messages/page.tsx
  - [x] Add ping event listener (browser auto-responds with pong)
  - [x] Improved connection handling

- [ ] Update frontend/src/app/dashboard/page.tsx
  - [ ] Add ping event listener
  - [ ] Improved connection handling

- [ ] Update frontend/src/app/notifications/page.tsx
  - [ ] Add ping event listener
  - [ ] Improved connection handling

- [ ] Update frontend/src/app/profile/[id]/page.tsx
  - [ ] Add ping event listener
  - [ ] Improved connection handling

- [ ] Check frontend/src/app/groups/[id]/page.tsx
  - [ ] Add ping event listener if WebSocket is used

### Testing
- [ ] Test connection stability over 15+ minutes
- [ ] Verify automatic reconnection
- [ ] Check ping/pong activity in console
- [ ] Test with multiple browser tabs

## Expected Outcome
WebSocket connections should remain stable indefinitely with proper keepalive mechanism.
