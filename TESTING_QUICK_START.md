# Quick Start Guide for Testing Social Network

## Prerequisites

1. **Install Required Software**:
   - Docker Desktop (recommended) OR
   - Go 1.21+ and Node.js 18+
   - 2 different web browsers (Chrome, Firefox, Edge, etc.)

2. **Clone/Navigate to Project**:
   ```bash
   cd c:/Users/amani/social-network
   ```

---

## Option 1: Run with Docker (Recommended)

### Start the Application
```bash
# Build and start both backend and frontend
docker-compose up --build

# Or run in detached mode (background)
docker-compose up -d --build
```

### Verify Services are Running
- Backend: http://localhost:8080
- Frontend: http://localhost:3000

### Stop the Application
```bash
docker-compose down
```

---

## Option 2: Run Manually (Without Docker)

### Terminal 1 - Start Backend
```bash
cd backend
go run server.go
```

Backend should start on: http://localhost:8080

### Terminal 2 - Start Frontend
```bash
cd frontend
npm install
npm run dev
```

Frontend should start on: http://localhost:3000

---

## Quick Test Setup

### 1. Create Test Users

Open http://localhost:3000 in your browser and register these test users:

**User 1 (Private Profile)**:
- Email: alice@test.com
- Password: Test123!
- First Name: Alice
- Last Name: Smith
- Date of Birth: 1990-01-01
- Profile: Private

**User 2 (Public Profile)**:
- Email: bob@test.com
- Password: Test123!
- First Name: Bob
- Last Name: Johnson
- Date of Birth: 1992-05-15
- Profile: Public

**User 3 (Public Profile)**:
- Email: charlie@test.com
- Password: Test123!
- First Name: Charlie
- Last Name: Brown
- Date of Birth: 1988-12-20
- Profile: Public

### 2. Open Multiple Browsers

**Browser 1 (Chrome)**:
- Navigate to: http://localhost:3000
- Login as: alice@test.com

**Browser 2 (Firefox)**:
- Navigate to: http://localhost:3000
- Login as: bob@test.com

**Browser 3 (Edge - Optional)**:
- Navigate to: http://localhost:3000
- Login as: charlie@test.com

---

## Quick Test Scenarios

### Scenario 1: Test Follow System (5 minutes)

1. **Browser 2 (Bob)**: Search for Alice and try to follow her
   - Expected: Shows "Request Pending" (Alice is private)

2. **Browser 1 (Alice)**: Check notifications/pending requests
   - Expected: See follow request from Bob
   - Action: Accept the request

3. **Browser 2 (Bob)**: Refresh and check Alice's profile
   - Expected: Now following Alice, can see her posts

4. **Browser 2 (Bob)**: Search for Charlie and follow him
   - Expected: Immediately follows (Charlie is public)

### Scenario 2: Test Posts (5 minutes)

1. **Browser 1 (Alice)**: Create a post
   - Content: "Hello from Alice!"
   - Privacy: "Followers Only"
   - Action: Submit

2. **Browser 2 (Bob)**: Check feed
   - Expected: See Alice's post (Bob follows Alice)

3. **Browser 3 (Charlie)**: Check feed
   - Expected: Don't see Alice's post (Charlie doesn't follow Alice)

4. **Browser 2 (Bob)**: Comment on Alice's post
   - Comment: "Great post, Alice!"

5. **Browser 1 (Alice)**: Check notifications
   - Expected: See notification about Bob's comment

### Scenario 3: Test Groups (10 minutes)

1. **Browser 1 (Alice)**: Create a group
   - Name: "Test Group"
   - Description: "A group for testing"
   - Action: Invite Bob

2. **Browser 2 (Bob)**: Check notifications
   - Expected: See group invitation
   - Action: Accept invitation

3. **Browser 3 (Charlie)**: Find the group and request to join
   - Action: Click "Request to Join"

4. **Browser 1 (Alice)**: Check group join requests
   - Expected: See Charlie's request
   - Action: Accept request

5. **Browser 2 (Bob)**: Create a post in the group
   - Content: "Hello group!"

6. **Browser 1 (Alice)**: Create an event in the group
   - Title: "Group Meetup"
   - Description: "Let's meet!"
   - Date: Tomorrow
   - Time: 2:00 PM

7. **Browser 2 (Bob)**: View event and respond
   - Action: Click "Going"

8. **Browser 1 (Alice)**: Check event responses
   - Expected: See Bob marked as "Going"

---

## Common Issues & Solutions

### Issue: Backend won't start
**Solution**: 
- Check if port 8080 is already in use
- Ensure Go is installed: `go version`
- Check for errors in terminal output

### Issue: Frontend won't start
**Solution**:
- Check if port 3000 is already in use
- Run `npm install` in frontend directory
- Check for errors in terminal output

### Issue: Can't connect to backend
**Solution**:
- Verify backend is running on port 8080
- Check CORS settings in backend
- Clear browser cache and cookies

### Issue: Database errors
**Solution**:
- Delete `backend/social_network.db` and restart
- Migrations will recreate the database

### Issue: Images not uploading
**Solution**:
- Check `backend/uploads` directory exists
- Verify file permissions
- Check file size limits

---

## Testing Checklist

Use the provided `TESTING_CHECKLIST.md` to track your progress:

```bash
# Open the checklist
notepad TESTING_CHECKLIST.md
```

For detailed test procedures, see `TESTING_PLAN.md`:

```bash
# Open the detailed plan
notepad TESTING_PLAN.md
```

---

## Reporting Issues

When you find a bug, document:

1. **What you were doing**: Exact steps
2. **What you expected**: Expected behavior
3. **What happened**: Actual behavior
4. **Screenshots**: Visual evidence
5. **Browser**: Which browser and version
6. **Console errors**: Any errors in browser console (F12)

---

## Tips for Effective Testing

1. **Clear cookies** between test runs to simulate fresh users
2. **Use incognito/private windows** for additional test users
3. **Check browser console** (F12) for JavaScript errors
4. **Test on different screen sizes** (desktop, tablet, mobile)
5. **Try edge cases**: empty inputs, very long text, special characters
6. **Test error handling**: wrong passwords, invalid data, etc.

---

## Next Steps

1. ✅ Start the application (Docker or manual)
2. ✅ Create test users
3. ✅ Open multiple browsers
4. ✅ Run quick test scenarios
5. ✅ Use TESTING_CHECKLIST.md to track progress
6. ✅ Refer to TESTING_PLAN.md for detailed tests
7. ✅ Document any bugs found

---

## Need Help?

- Check `TESTING_PLAN.md` for detailed test procedures
- Review `TODO.md` to see what features are implemented
- Check backend logs for API errors
- Check browser console for frontend errors

Happy Testing! 🚀
