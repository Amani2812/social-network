package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID             int       `json:"id"`
	Email          string    `json:"email"`
	PasswordHash   string    `json:"-"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	DateOfBirth    time.Time `json:"date_of_birth"`
	AvatarPath     *string   `json:"avatar_path,omitempty"`
	Nickname       *string   `json:"nickname,omitempty"`
	AboutMe        *string   `json:"about_me,omitempty"`
	IsPublic       bool      `json:"is_public"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user *User) error {
	query := `
		INSERT INTO users (email, password_hash, first_name, last_name, date_of_birth, avatar_path, nickname, about_me, is_public, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := r.DB.Exec(query, user.Email, user.PasswordHash, user.FirstName, user.LastName, user.DateOfBirth, user.AvatarPath, user.Nickname, user.AboutMe, user.IsPublic, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(id)
	return nil
}

func (r *UserRepository) GetByID(id int) (*User, error) {
	query := `SELECT id, email, password_hash, first_name, last_name, date_of_birth, avatar_path, nickname, about_me, is_public, created_at, updated_at FROM users WHERE id = ?`
	row := r.DB.QueryRow(query, id)

	user := &User{}
	var avatarPath, nickname, aboutMe sql.NullString
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName, &user.DateOfBirth, &avatarPath, &nickname, &aboutMe, &user.IsPublic, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if avatarPath.Valid {
		user.AvatarPath = &avatarPath.String
	}
	if nickname.Valid {
		user.Nickname = &nickname.String
	}
	if aboutMe.Valid {
		user.AboutMe = &aboutMe.String
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (*User, error) {
	query := `SELECT id, email, password_hash, first_name, last_name, date_of_birth, avatar_path, nickname, about_me, is_public, created_at, updated_at FROM users WHERE email = ?`
	row := r.DB.QueryRow(query, email)

	user := &User{}
	var avatarPath, nickname, aboutMe sql.NullString
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName, &user.DateOfBirth, &avatarPath, &nickname, &aboutMe, &user.IsPublic, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if avatarPath.Valid {
		user.AvatarPath = &avatarPath.String
	}
	if nickname.Valid {
		user.Nickname = &nickname.String
	}
	if aboutMe.Valid {
		user.AboutMe = &aboutMe.String
	}

	return user, nil
}

func (r *UserRepository) Update(user *User) error {
	query := `
		UPDATE users
		SET email = ?, first_name = ?, last_name = ?, date_of_birth = ?, avatar_path = ?, nickname = ?, about_me = ?, is_public = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := r.DB.Exec(query, user.Email, user.FirstName, user.LastName, user.DateOfBirth, user.AvatarPath, user.Nickname, user.AboutMe, user.IsPublic, user.UpdatedAt, user.ID)
	return err
}

func (r *UserRepository) SearchByName(query string) ([]*User, error) {
	searchQuery := "%" + query + "%"
	queryStr := `SELECT id, email, first_name, last_name, avatar_path, nickname, is_public FROM users WHERE (first_name LIKE ? OR last_name LIKE ? OR nickname LIKE ?) AND is_public = 1`
	rows, err := r.DB.Query(queryStr, searchQuery, searchQuery, searchQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		var avatarPath, nickname sql.NullString
		err := rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &avatarPath, &nickname, &user.IsPublic)
		if err != nil {
			return nil, err
		}
		if avatarPath.Valid {
			user.AvatarPath = &avatarPath.String
		}
		if nickname.Valid {
			user.Nickname = &nickname.String
		}
		users = append(users, user)
	}
	return users, nil
}
