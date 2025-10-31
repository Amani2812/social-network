package models

import (
	"database/sql"
	"time"
)

type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	ImagePath *string   `json:"image_path,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CommentRepository struct {
	DB *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{DB: db}
}

func (r *CommentRepository) Create(comment *Comment) error {
	query := `
		INSERT INTO comments (post_id, user_id, content, image_path, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := r.DB.Exec(query, comment.PostID, comment.UserID, comment.Content, comment.ImagePath, comment.CreatedAt, comment.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	comment.ID = int(id)
	return nil
}

func (r *CommentRepository) GetByPostID(postID int) ([]*Comment, error) {
	query := `SELECT id, post_id, user_id, content, image_path, created_at, updated_at FROM comments WHERE post_id = ? ORDER BY created_at ASC`
	rows, err := r.DB.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*Comment
	for rows.Next() {
		comment := &Comment{}
		var imagePath sql.NullString
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &imagePath, &comment.CreatedAt, &comment.UpdatedAt)
		if err != nil {
			return nil, err
		}
		if imagePath.Valid {
			comment.ImagePath = &imagePath.String
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (r *CommentRepository) GetByID(id int) (*Comment, error) {
	query := `SELECT id, post_id, user_id, content, image_path, created_at, updated_at FROM comments WHERE id = ?`
	row := r.DB.QueryRow(query, id)

	comment := &Comment{}
	var imagePath sql.NullString
	err := row.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &imagePath, &comment.CreatedAt, &comment.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if imagePath.Valid {
		comment.ImagePath = &imagePath.String
	}

	return comment, nil
}

func (r *CommentRepository) Update(comment *Comment) error {
	query := `UPDATE comments SET content = ?, image_path = ?, updated_at = ? WHERE id = ?`
	_, err := r.DB.Exec(query, comment.Content, comment.ImagePath, comment.UpdatedAt, comment.ID)
	return err
}

func (r *CommentRepository) Delete(id int) error {
	query := `DELETE FROM comments WHERE id = ?`
	_, err := r.DB.Exec(query, id)
	return err
}
