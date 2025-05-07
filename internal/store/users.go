package store

import (
	"context"
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
	CreatedAt string   `json:"created_at"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.text = &text
	p.hash = hash
	return nil
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) DeleteAll(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, "TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	return err
}
func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users(username,email,password)
		VALUES ($1,$2,$3)
		RETURNING id,email,created_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	err := s.db.QueryRowContext(ctx, query,
		user.Username,
		user.Email,
		user.Password,
	).Scan(&user.ID, &user.Email, &user.CreatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) GetById(ctx context.Context, userId int64) (*User, error) {
	query := `
		SELECT id,username,email,created_at FROM users WHERE id = $1;
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}
	err := s.db.QueryRowContext(ctx, query,
		userId,
	).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrorNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}
