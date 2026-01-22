package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	DateOfBirth  string    `json:"date_of_birth"`
	AvatarPath   *string   `json:"avatar_path,omitempty"`
	Nickname     *string   `json:"nickname,omitempty"`
	AboutMe      *string   `json:"about_me,omitempty"`
	IsPrivate    bool      `json:"is_private"`
	CreatedAt    time.Time `json:"created_at"`
}

// Session represents a user session
type Session struct {
	ID        string    `json:"id"`
	UserID    int       `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// Follow represents a follow relationship
type Follow struct {
	ID          int       `json:"id"`
	FollowerID  int       `json:"follower_id"`
	FollowingID int       `json:"following_id"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

// Post represents a user post
type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	ImagePath *string   `json:"image_path,omitempty"`
	Privacy   string    `json:"privacy"`
	CreatedAt time.Time `json:"created_at"`
	User      *User     `json:"user,omitempty"`
}

// Comment represents a comment on a post
type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	ImagePath *string   `json:"image_path,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	User      *User     `json:"user,omitempty"`
}

// Group represents a group
type Group struct {
	ID          int       `json:"id"`
	CreatorID   int       `json:"creator_id"`
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// GroupMember represents a group membership
type GroupMember struct {
	ID        int       `json:"id"`
	GroupID   int       `json:"group_id"`
	UserID    int       `json:"user_id"`
	Status    string    `json:"status"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	User      *User     `json:"user,omitempty"`
}

// GroupPost represents a post in a group
type GroupPost struct {
	ID        int       `json:"id"`
	GroupID   int       `json:"group_id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	ImagePath *string   `json:"image_path,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	User      *User     `json:"user,omitempty"`
}

// GroupEvent represents an event in a group
type GroupEvent struct {
	ID             int       `json:"id"`
	GroupID        int       `json:"group_id"`
	CreatorID      int       `json:"creator_id"`
	Title          string    `json:"title"`
	Description    *string   `json:"description,omitempty"`
	EventTime      time.Time `json:"event_time"`
	CreatedAt      time.Time `json:"created_at"`
	UserResponse   *string   `json:"user_response,omitempty"`
	GoingCount     int       `json:"going_count"`
	NotGoingCount  int       `json:"not_going_count"`
}

// EventResponse represents a user's response to an event
type EventResponse struct {
	ID        int       `json:"id"`
	EventID   int       `json:"event_id"`
	UserID    int       `json:"user_id"`
	Response  string    `json:"response"`
	CreatedAt time.Time `json:"created_at"`
}

// Message represents a chat message
type Message struct {
	ID         int       `json:"id"`
	SenderID   int       `json:"sender_id"`
	ReceiverID *int      `json:"receiver_id,omitempty"`
	GroupID    *int      `json:"group_id,omitempty"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	Sender     *User     `json:"sender,omitempty"`
}

// Notification represents a notification
type Notification struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	RelatedID *int      `json:"related_id,omitempty"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

// Repository contains all data access methods
type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// User methods
func (r *Repository) CreateUser(email, password, firstName, lastName, dob string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	result, err := r.db.Exec(
		"INSERT INTO users (email, password_hash, first_name, last_name, date_of_birth) VALUES (?, ?, ?, ?, ?)",
		email, string(hash), firstName, lastName, dob,
	)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return r.GetUserByID(int(id))
}

func (r *Repository) GetUserByEmail(email string) (*User, error) {
	user := &User{}
	err := r.db.QueryRow(
		"SELECT id, email, password_hash, first_name, last_name, date_of_birth, avatar_path, nickname, about_me, is_private, created_at FROM users WHERE email = ?",
		email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName, &user.DateOfBirth, &user.AvatarPath, &user.Nickname, &user.AboutMe, &user.IsPrivate, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) GetUserByID(id int) (*User, error) {
	user := &User{}
	err := r.db.QueryRow(
		"SELECT id, email, password_hash, first_name, last_name, date_of_birth, avatar_path, nickname, about_me, is_private, created_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName, &user.DateOfBirth, &user.AvatarPath, &user.Nickname, &user.AboutMe, &user.IsPrivate, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) UpdateUserProfile(userID int, avatarPath, nickname, aboutMe *string, isPrivate bool) error {
	_, err := r.db.Exec(
		"UPDATE users SET avatar_path = ?, nickname = ?, about_me = ?, is_private = ? WHERE id = ?",
		avatarPath, nickname, aboutMe, isPrivate, userID,
	)
	return err
}

func (r *Repository) SearchUsers(query string) ([]*User, error) {
	rows, err := r.db.Query(
		"SELECT id, email, first_name, last_name, avatar_path, nickname FROM users WHERE first_name LIKE ? OR last_name LIKE ? OR nickname LIKE ? LIMIT 20",
		"%"+query+"%", "%"+query+"%", "%"+query+"%",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		if err := rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.AvatarPath, &user.Nickname); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// Session methods
func (r *Repository) CreateSession(userID int) (*Session, error) {
	sessionID := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)

	_, err := r.db.Exec(
		"INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)",
		sessionID, userID, expiresAt,
	)
	if err != nil {
		return nil, err
	}

	return &Session{
		ID:        sessionID,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}, nil
}

func (r *Repository) GetSession(sessionID string) (*Session, error) {
	session := &Session{}
	err := r.db.QueryRow(
		"SELECT id, user_id, expires_at, created_at FROM sessions WHERE id = ? AND expires_at > ?",
		sessionID, time.Now(),
	).Scan(&session.ID, &session.UserID, &session.ExpiresAt, &session.CreatedAt)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (r *Repository) DeleteSession(sessionID string) error {
	_, err := r.db.Exec("DELETE FROM sessions WHERE id = ?", sessionID)
	return err
}

// Follow methods
func (r *Repository) CreateFollowRequest(followerID, followingID int) error {
	user, err := r.GetUserByID(followingID)
	if err != nil {
		return err
	}

	status := "accepted"
	if user.IsPrivate {
		status = "pending"
	}

	_, err = r.db.Exec(
		"INSERT INTO follows (follower_id, following_id, status) VALUES (?, ?, ?)",
		followerID, followingID, status,
	)
	return err
}

func (r *Repository) UpdateFollowStatus(followerID, followingID int, status string) error {
	_, err := r.db.Exec(
		"UPDATE follows SET status = ? WHERE follower_id = ? AND following_id = ?",
		status, followerID, followingID,
	)
	return err
}

func (r *Repository) DeleteFollow(followerID, followingID int) error {
	_, err := r.db.Exec(
		"DELETE FROM follows WHERE follower_id = ? AND following_id = ?",
		followerID, followingID,
	)
	return err
}

func (r *Repository) GetFollowers(userID int) ([]*User, error) {
	rows, err := r.db.Query(
		`SELECT u.id, u.email, u.first_name, u.last_name, u.avatar_path, u.nickname 
		FROM users u 
		JOIN follows f ON u.id = f.follower_id 
		WHERE f.following_id = ? AND f.status = 'accepted'`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		if err := rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.AvatarPath, &user.Nickname); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *Repository) GetFollowing(userID int) ([]*User, error) {
	rows, err := r.db.Query(
		`SELECT u.id, u.email, u.first_name, u.last_name, u.avatar_path, u.nickname 
		FROM users u 
		JOIN follows f ON u.id = f.following_id 
		WHERE f.follower_id = ? AND f.status = 'accepted'`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		if err := rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.AvatarPath, &user.Nickname); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *Repository) GetPendingFollowRequests(userID int) ([]*User, error) {
	rows, err := r.db.Query(
		`SELECT u.id, u.email, u.first_name, u.last_name, u.avatar_path, u.nickname 
		FROM users u 
		JOIN follows f ON u.id = f.follower_id 
		WHERE f.following_id = ? AND f.status = 'pending'`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		if err := rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.AvatarPath, &user.Nickname); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *Repository) IsFollowing(followerID, followingID int) (bool, error) {
	var count int
	err := r.db.QueryRow(
		"SELECT COUNT(*) FROM follows WHERE follower_id = ? AND following_id = ? AND status = 'accepted'",
		followerID, followingID,
	).Scan(&count)
	return count > 0, err
}

func (r *Repository) GetFollowRelationship(followerID, followingID int) (*Follow, error) {
	follow := &Follow{}
	err := r.db.QueryRow(
		"SELECT id, follower_id, following_id, status, created_at FROM follows WHERE follower_id = ? AND following_id = ?",
		followerID, followingID,
	).Scan(&follow.ID, &follow.FollowerID, &follow.FollowingID, &follow.Status, &follow.CreatedAt)
	if err != nil {
		return nil, err
	}
	return follow, nil
}

// Post methods
func (r *Repository) CreatePost(userID int, content, privacy string, imagePath *string) (*Post, error) {
	result, err := r.db.Exec(
		"INSERT INTO posts (user_id, content, privacy, image_path) VALUES (?, ?, ?, ?)",
		userID, content, privacy, imagePath,
	)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return r.GetPost(int(id))
}

func (r *Repository) GetPost(postID int) (*Post, error) {
	post := &Post{}
	err := r.db.QueryRow(
		"SELECT id, user_id, content, image_path, privacy, created_at FROM posts WHERE id = ?",
		postID,
	).Scan(&post.ID, &post.UserID, &post.Content, &post.ImagePath, &post.Privacy, &post.CreatedAt)
	if err != nil {
		return nil, err
	}

	post.User, _ = r.GetUserByID(post.UserID)
	return post, nil
}

func (r *Repository) GetUserPosts(userID, viewerID int) ([]*Post, error) {
	query := `SELECT id, user_id, content, image_path, privacy, created_at FROM posts WHERE user_id = ?`
	
	if userID != viewerID {
		isFollowing, _ := r.IsFollowing(viewerID, userID)
		if isFollowing {
			query += " AND privacy IN ('public', 'almost_private')"
		} else {
			query += " AND privacy = 'public'"
		}
	}
	
	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		post := &Post{}
		if err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.ImagePath, &post.Privacy, &post.CreatedAt); err != nil {
			return nil, err
		}
		post.User, _ = r.GetUserByID(post.UserID)
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *Repository) GetFeed(userID int) ([]*Post, error) {
	rows, err := r.db.Query(
		`SELECT p.id, p.user_id, p.content, p.image_path, p.privacy, p.created_at 
		FROM posts p
		WHERE (p.user_id = ? OR 
			   (p.privacy = 'public') OR 
			   (p.privacy = 'almost_private' AND EXISTS (
				   SELECT 1 FROM follows f 
				   WHERE f.follower_id = ? AND f.following_id = p.user_id AND f.status = 'accepted'
			   )))
		ORDER BY p.created_at DESC
		LIMIT 50`,
		userID, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		post := &Post{}
		if err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.ImagePath, &post.Privacy, &post.CreatedAt); err != nil {
			return nil, err
		}
		post.User, _ = r.GetUserByID(post.UserID)
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *Repository) UpdatePost(postID, userID int, content string) error {
	_, err := r.db.Exec(
		"UPDATE posts SET content = ? WHERE id = ? AND user_id = ?",
		content, postID, userID,
	)
	return err
}

func (r *Repository) DeletePost(postID, userID int) error {
	_, err := r.db.Exec("DELETE FROM posts WHERE id = ? AND user_id = ?", postID, userID)
	return err
}

// Comment methods
func (r *Repository) CreateComment(postID, userID int, content string, imagePath *string) (*Comment, error) {
	result, err := r.db.Exec(
		"INSERT INTO comments (post_id, user_id, content, image_path) VALUES (?, ?, ?, ?)",
		postID, userID, content, imagePath,
	)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return r.GetComment(int(id))
}

func (r *Repository) GetComment(commentID int) (*Comment, error) {
	comment := &Comment{}
	err := r.db.QueryRow(
		"SELECT id, post_id, user_id, content, image_path, created_at FROM comments WHERE id = ?",
		commentID,
	).Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.ImagePath, &comment.CreatedAt)
	if err != nil {
		return nil, err
	}

	comment.User, _ = r.GetUserByID(comment.UserID)
	return comment, nil
}

func (r *Repository) GetPostComments(postID int) ([]*Comment, error) {
	rows, err := r.db.Query(
		"SELECT id, post_id, user_id, content, image_path, created_at FROM comments WHERE post_id = ? ORDER BY created_at ASC",
		postID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*Comment
	for rows.Next() {
		comment := &Comment{}
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.ImagePath, &comment.CreatedAt); err != nil {
			return nil, err
		}
		comment.User, _ = r.GetUserByID(comment.UserID)
		comments = append(comments, comment)
	}
	return comments, nil
}

// Group methods
func (r *Repository) CreateGroup(creatorID int, title string, description *string) (*Group, error) {
	result, err := r.db.Exec(
		"INSERT INTO groups (creator_id, title, description) VALUES (?, ?, ?)",
		creatorID, title, description,
	)
	if err != nil {
		return nil, err
	}

	groupID, _ := result.LastInsertId()
	
	// Add creator as admin
	_, err = r.db.Exec(
		"INSERT INTO group_members (group_id, user_id, status, role) VALUES (?, ?, 'accepted', 'admin')",
		groupID, creatorID,
	)
	if err != nil {
		return nil, err
	}

	return r.GetGroup(int(groupID))
}

func (r *Repository) GetGroup(groupID int) (*Group, error) {
	group := &Group{}
	err := r.db.QueryRow(
		"SELECT id, creator_id, title, description, created_at FROM groups WHERE id = ?",
		groupID,
	).Scan(&group.ID, &group.CreatorID, &group.Title, &group.Description, &group.CreatedAt)
	return group, err
}

func (r *Repository) GetUserGroups(userID int) ([]*Group, error) {
	rows, err := r.db.Query(
		`SELECT g.id, g.creator_id, g.title, g.description, g.created_at 
		FROM groups g
		JOIN group_members gm ON g.id = gm.group_id
		WHERE gm.user_id = ? AND gm.status = 'accepted'
		ORDER BY g.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*Group
	for rows.Next() {
		group := &Group{}
		if err := rows.Scan(&group.ID, &group.CreatorID, &group.Title, &group.Description, &group.CreatedAt); err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func (r *Repository) GetAllGroups() ([]*Group, error) {
	rows, err := r.db.Query(
		`SELECT id, creator_id, title, description, created_at 
		FROM groups 
		ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*Group
	for rows.Next() {
		group := &Group{}
		if err := rows.Scan(&group.ID, &group.CreatorID, &group.Title, &group.Description, &group.CreatedAt); err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func (r *Repository) InviteToGroup(groupID, userID, inviterID int) error {
	// Check if inviter is a member (admin or regular member)
	var count int
	err := r.db.QueryRow(
		"SELECT COUNT(*) FROM group_members WHERE group_id = ? AND user_id = ? AND status = 'accepted'",
		groupID, inviterID,
	).Scan(&count)
	if err != nil || count == 0 {
		return errors.New("only group members can invite users")
	}

	_, err = r.db.Exec(
		"INSERT INTO group_members (group_id, user_id, status, role) VALUES (?, ?, 'pending', 'member')",
		groupID, userID,
	)
	return err
}

func (r *Repository) RequestToJoinGroup(groupID, userID int) error {
	// Check if user is already a member or has pending request
	var count int
	err := r.db.QueryRow(
		"SELECT COUNT(*) FROM group_members WHERE group_id = ? AND user_id = ?",
		groupID, userID,
	).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("user already has a relationship with this group")
	}

	_, err = r.db.Exec(
		"INSERT INTO group_members (group_id, user_id, status, role) VALUES (?, ?, 'pending', 'member')",
		groupID, userID,
	)
	return err
}

func (r *Repository) GetGroupJoinRequests(groupID int) ([]*GroupMember, error) {
	rows, err := r.db.Query(
		`SELECT gm.id, gm.group_id, gm.user_id, gm.status, gm.role, gm.created_at,
		u.id, u.email, u.first_name, u.last_name, u.avatar_path, u.nickname
		FROM group_members gm
		JOIN users u ON gm.user_id = u.id
		WHERE gm.group_id = ? AND gm.status = 'pending'`,
		groupID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*GroupMember
	for rows.Next() {
		member := &GroupMember{}
		user := &User{}
		if err := rows.Scan(
			&member.ID, &member.GroupID, &member.UserID, &member.Status, &member.Role, &member.CreatedAt,
			&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.AvatarPath, &user.Nickname,
		); err != nil {
			return nil, err
		}
		member.User = user
		requests = append(requests, member)
	}
	return requests, nil
}

func (r *Repository) RespondToJoinRequest(groupID, userID, adminID int, accept bool) error {
	// Check if responder is admin
	var role string
	err := r.db.QueryRow(
		"SELECT role FROM group_members WHERE group_id = ? AND user_id = ? AND status = 'accepted'",
		groupID, adminID,
	).Scan(&role)
	if err != nil || role != "admin" {
		return errors.New("only admins can respond to join requests")
	}

	status := "rejected"
	if accept {
		status = "accepted"
	}

	_, err = r.db.Exec(
		"UPDATE group_members SET status = ? WHERE group_id = ? AND user_id = ?",
		status, groupID, userID,
	)
	return err
}

func (r *Repository) RespondToGroupInvite(groupID, userID int, accept bool) error {
	status := "rejected"
	if accept {
		status = "accepted"
	}

	_, err := r.db.Exec(
		"UPDATE group_members SET status = ? WHERE group_id = ? AND user_id = ?",
		status, groupID, userID,
	)
	return err
}

func (r *Repository) CreateGroupPost(groupID, userID int, content string, imagePath *string) (*GroupPost, error) {
	// Check if user is member
	var count int
	err := r.db.QueryRow(
		"SELECT COUNT(*) FROM group_members WHERE group_id = ? AND user_id = ? AND status = 'accepted'",
		groupID, userID,
	).Scan(&count)
	if err != nil || count == 0 {
		return nil, errors.New("user is not a member of this group")
	}

	result, err := r.db.Exec(
		"INSERT INTO group_posts (group_id, user_id, content, image_path) VALUES (?, ?, ?, ?)",
		groupID, userID, content, imagePath,
	)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return r.GetGroupPost(int(id))
}

func (r *Repository) GetGroupPost(postID int) (*GroupPost, error) {
	post := &GroupPost{}
	err := r.db.QueryRow(
		"SELECT id, group_id, user_id, content, image_path, created_at FROM group_posts WHERE id = ?",
		postID,
	).Scan(&post.ID, &post.GroupID, &post.UserID, &post.Content, &post.ImagePath, &post.CreatedAt)
	if err != nil {
		return nil, err
	}

	post.User, _ = r.GetUserByID(post.UserID)
	return post, nil
}

func (r *Repository) GetGroupPosts(groupID, userID int) ([]*GroupPost, error) {
	// Check if user is member
	var count int
	err := r.db.QueryRow(
		"SELECT COUNT(*) FROM group_members WHERE group_id = ? AND user_id = ? AND status = 'accepted'",
		groupID, userID,
	).Scan(&count)
	if err != nil || count == 0 {
		return nil, errors.New("user is not a member of this group")
	}

	rows, err := r.db.Query(
		"SELECT id, group_id, user_id, content, image_path, created_at FROM group_posts WHERE group_id = ? ORDER BY created_at DESC",
		groupID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*GroupPost
	for rows.Next() {
		post := &GroupPost{}
		if err := rows.Scan(&post.ID, &post.GroupID, &post.UserID, &post.Content, &post.ImagePath, &post.CreatedAt); err != nil {
			return nil, err
		}
		post.User, _ = r.GetUserByID(post.UserID)
		posts = append(posts, post)
	}
	return posts, nil
}

// Event methods
func (r *Repository) CreateEvent(groupID, creatorID int, title string, description *string, eventTime time.Time) (*GroupEvent, error) {
	result, err := r.db.Exec(
		"INSERT INTO group_events (group_id, creator_id, title, description, event_time) VALUES (?, ?, ?, ?, ?)",
		groupID, creatorID, title, description, eventTime,
	)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return r.GetEvent(int(id))
}

func (r *Repository) GetEvent(eventID int) (*GroupEvent, error) {
	event := &GroupEvent{}
	err := r.db.QueryRow(
		"SELECT id, group_id, creator_id, title, description, event_time, created_at FROM group_events WHERE id = ?",
		eventID,
	).Scan(&event.ID, &event.GroupID, &event.CreatorID, &event.Title, &event.Description, &event.EventTime, &event.CreatedAt)
	return event, err
}

func (r *Repository) GetGroupEvents(groupID int) ([]*GroupEvent, error) {
	rows, err := r.db.Query(
		"SELECT id, group_id, creator_id, title, description, event_time, created_at FROM group_events WHERE group_id = ? ORDER BY event_time ASC",
		groupID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*GroupEvent
	for rows.Next() {
		event := &GroupEvent{}
		if err := rows.Scan(&event.ID, &event.GroupID, &event.CreatorID, &event.Title, &event.Description, &event.EventTime, &event.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (r *Repository) GetGroupEventsWithUserResponse(groupID, userID int) ([]*GroupEvent, error) {
	rows, err := r.db.Query(
		`SELECT ge.id, ge.group_id, ge.creator_id, ge.title, ge.description, ge.event_time, ge.created_at,
		er.response,
		COALESCE((SELECT COUNT(*) FROM event_responses WHERE event_id = ge.id AND response = 'going'), 0) as going_count,
		COALESCE((SELECT COUNT(*) FROM event_responses WHERE event_id = ge.id AND response = 'not_going'), 0) as not_going_count
		FROM group_events ge
		LEFT JOIN event_responses er ON ge.id = er.event_id AND er.user_id = ?
		WHERE ge.group_id = ?
		ORDER BY ge.event_time ASC`,
		userID, groupID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*GroupEvent
	for rows.Next() {
		event := &GroupEvent{}
		var userResponse sql.NullString
		if err := rows.Scan(
			&event.ID, &event.GroupID, &event.CreatorID, &event.Title, &event.Description, &event.EventTime, &event.CreatedAt,
			&userResponse,
			&event.GoingCount,
			&event.NotGoingCount,
		); err != nil {
			return nil, err
		}
		if userResponse.Valid {
			event.UserResponse = &userResponse.String
		}
		events = append(events, event)
	}
	return events, nil
}

func (r *Repository) GetEventResponses(eventID int) ([]*EventResponse, error) {
	rows, err := r.db.Query(
		"SELECT id, event_id, user_id, response, created_at FROM event_responses WHERE event_id = ?",
		eventID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var responses []*EventResponse
	for rows.Next() {
		response := &EventResponse{}
		if err := rows.Scan(&response.ID, &response.EventID, &response.UserID, &response.Response, &response.CreatedAt); err != nil {
			return nil, err
		}
		responses = append(responses, response)
	}
	return responses, nil
}

func (r *Repository) RespondToEvent(eventID, userID int, response string) error {
	_, err := r.db.Exec(
		"INSERT OR REPLACE INTO event_responses (event_id, user_id, response) VALUES (?, ?, ?)",
		eventID, userID, response,
	)
	return err
}

// Message methods
func (r *Repository) CreateMessage(senderID int, receiverID, groupID *int, content string) (*Message, error) {
	result, err := r.db.Exec(
		"INSERT INTO messages (sender_id, receiver_id, group_id, content) VALUES (?, ?, ?, ?)",
		senderID, receiverID, groupID, content,
	)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return r.GetMessage(int(id))
}

func (r *Repository) GetMessage(messageID int) (*Message, error) {
	msg := &Message{}
	err := r.db.QueryRow(
		"SELECT id, sender_id, receiver_id, group_id, content, created_at FROM messages WHERE id = ?",
		messageID,
	).Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.GroupID, &msg.Content, &msg.CreatedAt)
	if err != nil {
		return nil, err
	}

	msg.Sender, _ = r.GetUserByID(msg.SenderID)
	return msg, nil
}

func (r *Repository) GetPrivateMessages(user1ID, user2ID int) ([]*Message, error) {
	rows, err := r.db.Query(
		`SELECT id, sender_id, receiver_id, group_id, content, created_at FROM messages 
		WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)
		ORDER BY created_at ASC`,
		user1ID, user2ID, user2ID, user1ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		msg := &Message{}
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.GroupID, &msg.Content, &msg.CreatedAt); err != nil {
			return nil, err
		}
		msg.Sender, _ = r.GetUserByID(msg.SenderID)
		messages = append(messages, msg)
	}
	return messages, nil
}

func (r *Repository) GetGroupMessages(groupID int) ([]*Message, error) {
	rows, err := r.db.Query(
		"SELECT id, sender_id, receiver_id, group_id, content, created_at FROM messages WHERE group_id = ? ORDER BY created_at ASC",
		groupID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		msg := &Message{}
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.GroupID, &msg.Content, &msg.CreatedAt); err != nil {
			return nil, err
		}
		msg.Sender, _ = r.GetUserByID(msg.SenderID)
		messages = append(messages, msg)
	}
	return messages, nil
}

func (r *Repository) GetGroupMembers(groupID int) ([]*GroupMember, error) {
	rows, err := r.db.Query(
		"SELECT id, group_id, user_id, status, role, created_at FROM group_members WHERE group_id = ? AND status = 'accepted'",
		groupID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*GroupMember
	for rows.Next() {
		member := &GroupMember{}
		if err := rows.Scan(&member.ID, &member.GroupID, &member.UserID, &member.Status, &member.Role, &member.CreatedAt); err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	return members, nil
}

func (r *Repository) GetGroupMembersWithDetails(groupID int) ([]*GroupMember, error) {
	rows, err := r.db.Query(
		`SELECT gm.id, gm.group_id, gm.user_id, gm.status, gm.role, gm.created_at,
		u.id, u.email, u.first_name, u.last_name, u.avatar_path, u.nickname
		FROM group_members gm
		JOIN users u ON gm.user_id = u.id
		WHERE gm.group_id = ? AND gm.status = 'accepted'`,
		groupID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*GroupMember
	for rows.Next() {
		member := &GroupMember{}
		user := &User{}
		if err := rows.Scan(
			&member.ID, &member.GroupID, &member.UserID, &member.Status, &member.Role, &member.CreatedAt,
			&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.AvatarPath, &user.Nickname,
		); err != nil {
			return nil, err
		}
		member.User = user
		members = append(members, member)
	}
	return members, nil
}

func (r *Repository) GetConversations(userID int) ([]*User, error) {
	rows, err := r.db.Query(
		`SELECT DISTINCT u.id, u.email, u.first_name, u.last_name, u.avatar_path, u.nickname
		FROM users u
		JOIN messages m ON (m.sender_id = u.id OR m.receiver_id = u.id)
		WHERE (m.sender_id = ? OR m.receiver_id = ?) AND u.id != ?
		ORDER BY (
			SELECT MAX(created_at) FROM messages 
			WHERE (sender_id = u.id AND receiver_id = ?) OR (sender_id = ? AND receiver_id = u.id)
		) DESC`,
		userID, userID, userID, userID, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		if err := rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.AvatarPath, &user.Nickname); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// Notification methods
func (r *Repository) CreateNotification(userID int, notifType, content string, relatedID *int) error {
	_, err := r.db.Exec(
		"INSERT INTO notifications (user_id, type, content, related_id) VALUES (?, ?, ?, ?)",
		userID, notifType, content, relatedID,
	)
	return err
}

func (r *Repository) GetUserNotifications(userID int) ([]*Notification, error) {
	rows, err := r.db.Query(
		"SELECT id, user_id, type, content, related_id, is_read, created_at FROM notifications WHERE user_id = ? ORDER BY created_at DESC LIMIT 50",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*Notification
	for rows.Next() {
		notif := &Notification{}
		if err := rows.Scan(&notif.ID, &notif.UserID, &notif.Type, &notif.Content, &notif.RelatedID, &notif.IsRead, &notif.CreatedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, notif)
	}
	return notifications, nil
}

func (r *Repository) MarkNotificationAsRead(notificationID, userID int) error {
	_, err := r.db.Exec(
		"UPDATE notifications SET is_read = 1 WHERE id = ? AND user_id = ?",
		notificationID, userID,
	)
	return err
}

func (r *Repository) GetUnreadNotificationCount(userID int) (int, error) {
	var count int
	err := r.db.QueryRow(
		"SELECT COUNT(*) FROM notifications WHERE user_id = ? AND is_read = 0",
		userID,
	).Scan(&count)
	return count, err
}
