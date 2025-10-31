package models

import (
	"database/sql"
	"time"
)

type Post struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Content     string    `json:"content"`
	ImagePath   *string   `json:"image_path,omitempty"`
	Privacy     string    `json:"privacy"` // "public", "almost_private", "private"
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PostRepository struct {
	DB *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{DB: db}
}

func (r *PostRepository) Create(post *Post) error {
	query := `
		INSERT INTO posts (user_id, content, image_path, privacy, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := r.DB.Exec(query, post.UserID, post.Content, post.ImagePath, post.Privacy, post.CreatedAt, post.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	post.ID = int(id)
	return nil
}

func (r *PostRepository) GetByID(id int) (*Post, error) {
	query := `SELECT id, user_id, content, image_path, privacy, created_at, updated_at FROM posts WHERE id = ?`
	row := r.DB.QueryRow(query, id)

	post := &Post{}
	var imagePath sql.NullString
	err := row.Scan(&post.ID, &post.UserID, &post.Content, &imagePath, &post.Privacy, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if imagePath.Valid {
		post.ImagePath = &imagePath.String
	}

	return post, nil
}

func (r *PostRepository) GetByUserID(userID int, currentUserID int) ([]*Post, error) {
	// Check if current user can view posts
	var query string
	var args []interface{}

	if currentUserID == userID {
		// User can see all their own posts
		query = `SELECT id, user_id, content, image_path, privacy, created_at, updated_at FROM posts WHERE user_id = ? ORDER BY created_at DESC`
		args = []interface{}{userID}
	} else {
		// Check if current user follows the post owner
		followQuery := `SELECT COUNT(*) FROM follows WHERE follower_id = ? AND following_id = ? AND status = 'accepted'`
		var followCount int
		err := r.DB.QueryRow(followQuery, currentUserID, userID).Scan(&followCount)
		if err != nil {
			return nil, err
		}

		if followCount > 0 {
			// Following, can see public and almost_private posts
			query = `SELECT id, user_id, content, image_path, privacy, created_at, updated_at FROM posts WHERE user_id = ? AND privacy IN ('public', 'almost_private') ORDER BY created_at DESC`
			args = []interface{}{userID}
		} else {
			// Not following, can only see public posts
			query = `SELECT id, user_id, content, image_path, privacy, created_at, updated_at FROM posts WHERE user_id = ? AND privacy = 'public' ORDER BY created_at DESC`
			args = []interface{}{userID}
		}
	}

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		post := &Post{}
		var imagePath sql.NullString
		err := rows.Scan(&post.ID, &post.UserID, &post.Content, &imagePath, &post.Privacy, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}
		if imagePath.Valid {
			post.ImagePath = &imagePath.String
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostRepository) GetFeed(userID int) ([]*Post, error) {
	// Get posts from users that the current user follows, plus their own posts
	query := `
		SELECT p.id, p.user_id, p.content, p.image_path, p.privacy, p.created_at, p.updated_at
		FROM posts p
		INNER JOIN follows f ON p.user_id = f.following_id
		WHERE f.follower_id = ? AND f.status = 'accepted' AND p.privacy IN ('public', 'almost_private')
		UNION
		SELECT id, user_id, content, image_path, privacy, created_at, updated_at
		FROM posts
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT 50
	`

	rows, err := r.DB.Query(query, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		post := &Post{}
		var imagePath sql.NullString
		err := rows.Scan(&post.ID, &post.UserID, &post.Content, &imagePath, &post.Privacy, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}
		if imagePath.Valid {
			post.ImagePath = &imagePath.String
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostRepository) Update(post *Post) error {
	query := `UPDATE posts SET content = ?, image_path = ?, privacy = ?, updated_at = ? WHERE id = ?`
	_, err := r.DB.Exec(query, post.Content, post.ImagePath, post.Privacy, post.UpdatedAt, post.ID)
	return err
}

func (r *PostRepository) Delete(id int) error {
	query := `DELETE FROM posts WHERE id = ?`
	_, err := r.DB.Exec(query, id)
	return err
}
