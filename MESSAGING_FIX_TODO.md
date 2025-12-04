# Messaging Fix TODO

## Issues to Fix:
- [x] Identified root causes
- [x] Fix frontend WebSocket connection management
- [x] Add optimistic UI updates
- [x] Fix backend message echoing
- [ ] Test messaging between two accounts

## Root Causes:
1. Multiple WebSocket connections being created without proper cleanup
2. No optimistic UI updates (messages don't appear immediately)
3. Backend echoing messages back to sender causing confusion
4. WebSocket connection lifecycle not properly managed

## Files Edited:
1. ✅ frontend/src/app/messages/page.tsx - Fixed WebSocket management and added optimistic UI
2. ✅ backend/pkg/websocket/websocket.go - Removed echo to sender

## Changes Made:

### Frontend (messages/page.tsx):
1. Added `hasConnectedRef` to prevent multiple WebSocket connections
2. Improved WebSocket connection lifecycle management
3. Added optimistic UI updates - messages appear immediately when sent
4. Modified message handler to skip own messages (handled optimistically)
5. Improved duplicate detection (checks last 10 messages instead of 5)
6. Better WebSocket cleanup on unmount

### Backend (websocket.go):
1. Removed echo back to sender (`h.sendToUser(message.SenderID, message)`)
2. Only sends messages to receiver for real-time delivery
3. Sender sees their message via optimistic UI update

## How It Works Now:
1. User types a message and clicks Send
2. Message appears immediately in their chat (optimistic UI)
3. Message is sent via WebSocket to backend
4. Backend saves to database and sends to receiver only
5. Receiver gets message in real-time via WebSocket
6. No duplicate messages, no echo confusion

## Next Steps:
- Test with two browser windows/accounts
- Verify messages appear immediately when sent
- Verify receiver gets messages in real-time
- Test page refresh behavior
