package models

import (
	"database/sql"
	"time"
)

type Follow struct {
	ID          int       `json:"id"`
	FollowerID  int       `json:"follower_id"`
	FollowingID int       `json:"following_id"`
	Status      string    `json:"status"` // pending, accepted, rejected
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type FollowRepository struct {
	DB *sql.DB
}

func NewFollowRepository(db *sql.DB) *FollowRepository {
	return &FollowRepository{DB: db}
}

func (r *FollowRepository) Create(follow *Follow) error {
	query := `
		INSERT INTO follows (follower_id, following_id, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.DB.Exec(query, follow.FollowerID, follow.FollowingID, follow.Status, follow.CreatedAt, follow.UpdatedAt)
	return err
}

func (r *FollowRepository) GetByID(id int) (*Follow, error) {
	follow := &Follow{}
	query := `SELECT id, follower_id, following_id, status, created_at, updated_at FROM follows WHERE id = ?`
	err := r.DB.QueryRow(query, id).Scan(&follow.ID, &follow.FollowerID, &follow.FollowingID, &follow.Status, &follow.CreatedAt, &follow.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return follow, nil
}

func (r *FollowRepository) GetPendingRequests(userID int) ([]*Follow, error) {
	query := `SELECT id, follower_id, following_id, status, created_at, updated_at FROM follows WHERE following_id = ? AND status = 'pending'`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var follows []*Follow
	for rows.Next() {
		follow := &Follow{}
		err := rows.Scan(&follow.ID, &follow.FollowerID, &follow.FollowingID, &follow.Status, &follow.CreatedAt, &follow.UpdatedAt)
		if err != nil {
			return nil, err
		}
		follows = append(follows, follow)
	}
	return follows, nil
}

func (r *FollowRepository) GetFollowers(userID int) ([]*Follow, error) {
	query := `SELECT id, follower_id, following_id, status, created_at, updated_at FROM follows WHERE following_id = ? AND status = 'accepted'`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var follows []*Follow
	for rows.Next() {
		follow := &Follow{}
		err := rows.Scan(&follow.ID, &follow.FollowerID, &follow.FollowingID, &follow.Status, &follow.CreatedAt, &follow.UpdatedAt)
		if err != nil {
			return nil, err
		}
		follows = append(follows, follow)
	}
	return follows, nil
}

func (r *FollowRepository) GetFollowing(userID int) ([]*Follow, error) {
	query := `SELECT id, follower_id, following_id, status, created_at, updated_at FROM follows WHERE follower_id = ? AND status = 'accepted'`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var follows []*Follow
	for rows.Next() {
		follow := &Follow{}
		err := rows.Scan(&follow.ID, &follow.FollowerID, &follow.FollowingID, &follow.Status, &follow.CreatedAt, &follow.UpdatedAt)
		if err != nil {
			return nil, err
		}
		follows = append(follows, follow)
	}
	return follows, nil
}

func (r *FollowRepository) UpdateStatus(id int, status string) error {
	query := `UPDATE follows SET status = ?, updated_at = ? WHERE id = ?`
	_, err := r.DB.Exec(query, status, time.Now(), id)
	return err
}

func (r *FollowRepository) Delete(id int) error {
	query := `DELETE FROM follows WHERE id = ?`
	_, err := r.DB.Exec(query, id)
	return err
}

func (r *FollowRepository) IsFollowing(followerID, followingID int) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM follows WHERE follower_id = ? AND following_id = ? AND status = 'accepted'`
	err := r.DB.QueryRow(query, followerID, followingID).Scan(&count)
	return count > 0, err
}

func (r *FollowRepository) GetFollowRequest(followerID, followingID int) (*Follow, error) {
	follow := &Follow{}
	query := `SELECT id, follower_id, following_id, status, created_at, updated_at FROM follows WHERE follower_id = ? AND following_id = ?`
	err := r.DB.QueryRow(query, followerID, followingID).Scan(&follow.ID, &follow.FollowerID, &follow.FollowingID, &follow.Status, &follow.CreatedAt, &follow.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return follow, nil
}
