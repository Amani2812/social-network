# Bug Report Template

Use this template to report bugs found during testing. Copy and fill out for each bug.

---

## Bug Report #___

### Bug Information
- **Date Found**: _______________
- **Tester**: _______________
- **Severity**: ⬜ Critical | ⬜ High | ⬜ Medium | ⬜ Low
- **Status**: ⬜ New | ⬜ In Progress | ⬜ Fixed | ⬜ Closed

### Environment
- **Browser**: _______________
- **Browser Version**: _______________
- **Operating System**: _______________
- **Screen Resolution**: _______________

### Test Case
**Which test was being performed?**
(e.g., Test 1.1: Follow Request to Private User)

_______________________________________________

### Summary
**Brief description of the bug (one sentence)**

_______________________________________________

### Steps to Reproduce
1. _______________________________________________
2. _______________________________________________
3. _______________________________________________
4. _______________________________________________
5. _______________________________________________

### Expected Result
**What should have happened?**

_______________________________________________
_______________________________________________

### Actual Result
**What actually happened?**

_______________________________________________
_______________________________________________

### Screenshots
**Attach screenshots here or reference file names**

- Screenshot 1: _______________
- Screenshot 2: _______________
- Screenshot 3: _______________

### Console Errors
**Any errors from browser console (F12)?**

```
[Paste console errors here]
```

### Network Errors
**Any failed API requests (Network tab in F12)?**

```
[Paste network errors here]
```

### Additional Notes
**Any other relevant information**

_______________________________________________
_______________________________________________

### Workaround
**Is there a temporary workaround?**

⬜ Yes | ⬜ No

If yes, describe:
_______________________________________________

---

## Example Bug Reports

### Example 1: Critical Bug

## Bug Report #001

### Bug Information
- **Date Found**: 2024-01-15
- **Tester**: John Doe
- **Severity**: ☑️ Critical | ⬜ High | ⬜ Medium | ⬜ Low
- **Status**: ☑️ New | ⬜ In Progress | ⬜ Fixed | ⬜ Closed

### Environment
- **Browser**: Chrome
- **Browser Version**: 120.0.6099.109
- **Operating System**: Windows 11
- **Screen Resolution**: 1920x1080

### Test Case
Test 1.1: Follow Request to Private User

### Summary
Follow request to private user is not sent - button does nothing when clicked

### Steps to Reproduce
1. Login as User B (bob@test.com) in Chrome
2. Search for User A (alice@test.com) who has a private profile
3. Navigate to User A's profile page
4. Click the "Follow" button
5. Observe the result

### Expected Result
- Button should change to "Request Pending"
- Follow request should be sent to User A
- User A should receive a notification

### Actual Result
- Button does not change
- No request is sent
- User A receives no notification
- No error message is displayed

### Screenshots
- screenshot_follow_button_before.png
- screenshot_follow_button_after.png

### Console Errors
```
POST http://localhost:8080/api/follow/request 500 (Internal Server Error)
Error: Failed to send follow request
    at handleFollow (profile.tsx:45)
```

### Network Errors
```
Request URL: http://localhost:8080/api/follow/request
Request Method: POST
Status Code: 500 Internal Server Error
Response: {"error": "database connection failed"}
```

### Additional Notes
This is a critical bug as it prevents the core follow functionality from working. The issue appears to be on the backend - database connection error.

### Workaround
⬜ Yes | ☑️ No

---

### Example 2: Medium Bug

## Bug Report #002

### Bug Information
- **Date Found**: 2024-01-15
- **Tester**: Jane Smith
- **Severity**: ⬜ Critical | ⬜ High | ☑️ Medium | ⬜ Low
- **Status**: ☑️ New | ⬜ In Progress | ⬜ Fixed | ⬜ Closed

### Environment
- **Browser**: Firefox
- **Browser Version**: 121.0
- **Operating System**: Windows 11
- **Screen Resolution**: 1920x1080

### Test Case
Test 3.2: Create Post with Image (JPG/PNG)

### Summary
Uploaded images are not displayed in posts - only broken image icon shows

### Steps to Reproduce
1. Login as User A
2. Navigate to dashboard
3. Click "Create Post"
4. Enter text: "Check out this image!"
5. Click image upload button
6. Select a JPG image (test.jpg, 2MB)
7. Submit the post
8. View the post in feed

### Expected Result
- Image should be displayed in the post
- Image should be properly sized and formatted
- Image should be clickable to view full size

### Actual Result
- Post is created successfully
- Text content is displayed correctly
- Image shows as broken image icon (🖼️ with X)
- Browser shows 404 error when trying to load image

### Screenshots
- screenshot_post_with_broken_image.png
- screenshot_upload_dialog.png

### Console Errors
```
GET http://localhost:8080/uploads/undefined 404 (Not Found)
```

### Network Errors
```
Request URL: http://localhost:8080/uploads/undefined
Status Code: 404 Not Found
```

### Additional Notes
The image appears to upload successfully (no error during upload), but the image path stored in the database might be incorrect or the file is not being saved to the uploads directory.

### Workaround
☑️ Yes | ⬜ No

Workaround: Use external image hosting (e.g., Imgur) and paste the URL in the post text. Not ideal but allows testing other features.

---

## Severity Guidelines

### Critical
- Application crashes
- Data loss
- Security vulnerabilities
- Core features completely broken
- Blocks all testing

### High
- Major features don't work
- Significant impact on user experience
- Workaround is difficult or impossible
- Affects multiple users

### Medium
- Feature works but with issues
- Moderate impact on user experience
- Workaround exists
- Affects some users

### Low
- Minor visual issues
- Typos or formatting problems
- Minimal impact on functionality
- Easy workaround available

---

## Status Definitions

- **New**: Bug just discovered, not yet reviewed
- **In Progress**: Bug is being worked on
- **Fixed**: Bug has been fixed, awaiting verification
- **Closed**: Bug verified as fixed and closed

---

## Tips for Good Bug Reports

1. **Be Specific**: Provide exact steps, not general descriptions
2. **Be Objective**: Describe what happened, not what you think is wrong
3. **One Bug Per Report**: Don't combine multiple issues
4. **Include Evidence**: Screenshots, console logs, network errors
5. **Test Reproducibility**: Can you make it happen again?
6. **Check for Duplicates**: Has this bug been reported already?
7. **Provide Context**: What were you trying to accomplish?

---

## Quick Bug Report (For Minor Issues)

For minor issues, you can use this shortened format:

**Bug**: _______________________________________________
**Steps**: _______________________________________________
**Expected**: _______________________________________________
**Actual**: _______________________________________________
**Severity**: _______________

---

## Tracking Your Bug Reports

Keep a list of all bugs found:

| # | Summary | Severity | Status | Date |
|---|---------|----------|--------|------|
| 001 | Follow button not working | Critical | New | 2024-01-15 |
| 002 | Images not displaying | Medium | New | 2024-01-15 |
| 003 | Typo in welcome message | Low | Fixed | 2024-01-15 |
| 004 | ... | ... | ... | ... |

---

## Submitting Bug Reports

Once you've documented bugs:

1. Save this file with your bug reports
2. Share with the development team
3. Include all screenshots in a separate folder
4. Prioritize critical and high severity bugs
5. Follow up on bug fixes and verify they work

---

Good luck with bug hunting! 🐛🔍
