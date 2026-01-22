# Like/Dislike Feature Implementation Guide

## Current Status

The **frontend UI is complete** for like/dislike functionality on posts and comments. However, the **backend API endpoints are missing**, causing "Failed to fetch" errors.

## Errors Shown

```
Failed to fetch
- http://localhost:8080/api/posts/react
- http://localhost:8080/api/comments/react
```

## Frontend Implementation (✅ Complete)

### Files Modified:
1. `frontend/src/app/dashboard/page.tsx` - Post reactions
2. `frontend/src/components/Comments.tsx` - Comment reactions

### Features:
- 👍 Like and 👎 Dislike buttons
- Visual feedback (blue for like, red for dislike)
- Real-time count display
- Toggle functionality

## Backend Implementation (❌ Missing)

### Required API Endpoints:

#### 1. POST `/api/posts/react`
**Request Body:**
```json
{
  "post_id": 123,
  "reaction": "like" // or "dislike"
}
```

**Response:**
```json
{
  "success": true,
  "likes": 10,
  "dislikes": 2,
  "user_reaction": "like" // or "dislike" or null
}
```

**Logic:**
- If user hasn't reacted: Add the reaction
- If user has same reaction: Remove the reaction (toggle off)
- If user has different reaction: Change to new reaction
- Return updated counts and user's current reaction

#### 2. POST `/api/comments/react`
**Request Body:**
```json
{
  "comment_id": 456,
  "reaction": "like" // or "dislike"
}
```

**Response:**
```json
{
  "success": true,
  "likes": 5,
  "dislikes": 1,
  "user_reaction": "like" // or "dislike" or null
}
```

**Logic:** Same as posts

### Database Schema Needed:

#### Table: `post_reactions`
```sql
CREATE TABLE post_reactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    reaction TEXT NOT NULL CHECK(reaction IN ('like', 'dislike')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(post_id, user_id)
);
```

#### Table: `comment_reactions`
```sql
CREATE TABLE comment_reactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    comment_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    reaction TEXT NOT NULL CHECK(reaction IN ('like', 'dislike')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(comment_id, user_id)
);
```

### Backend Implementation Steps:

1. **Create Migration Files** for the new tables
2. **Update Models** (`backend/pkg/models/models.go`):
   - Add `PostReaction` and `CommentReaction` structs
   - Add methods: `CreateReaction`, `DeleteReaction`, `GetReactionCounts`, `GetUserReaction`

3. **Update Handlers** (`backend/pkg/handlers/handlers.go`):
   - Add `HandlePostReact` function
   - Add `HandleCommentReact` function

4. **Update Routes** (`backend/server.go`):
   - Add route: `POST /api/posts/react`
   - Add route: `POST /api/comments/react`

5. **Update Feed/Comments Queries**:
   - Modify `GetFeed` to include like/dislike counts and user's reaction
   - Modify `GetComments` to include like/dislike counts and user's reaction

## Temporary Solution

Until the backend is implemented, you can:

### Option 1: Comment Out the Functionality
Remove the like/dislike buttons temporarily by commenting out the relevant sections in:
- `frontend/src/app/dashboard/page.tsx` (lines ~565-590)
- `frontend/src/components/Comments.tsx` (lines ~275-300)

### Option 2: Mock the Response
Add temporary mock responses in the frontend to test the UI without backend.

## Testing After Implementation

Once backend is ready, test:
1. ✅ Like a post - count increases
2. ✅ Like again - count decreases (toggle off)
3. ✅ Like then dislike - like count decreases, dislike increases
4. ✅ Refresh page - reactions persist
5. ✅ Multiple users - each user's reactions are independent
6. ✅ Same tests for comments

## Priority

This feature is **non-critical** for core functionality. The social network works fine without it. Implement when ready or leave for future enhancement.
