package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
	IsActive  bool     `json:"is_active"`
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

func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
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
		user.Password.hash,
	).Scan(&user.ID, &user.Email, &user.CreatedAt)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
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

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id,username,email,password,created_at FROM users WHERE email = $1 AND is_active = true;
	`
	fmt.Println(query)
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}
	err := s.db.QueryRowContext(ctx, query,
		email,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password.hash, &user.CreatedAt)

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

func (s *UserStore) CreateAndInvite(ctx context.Context, user *User, token string, expiresIn time.Duration) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		// create the user
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}

		if err := s.createUserInvitation(ctx, tx, token, expiresIn, user.ID); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) Activate(ctx context.Context, token string) error {
	// because we need to perform multiple actions in the db, we're using a transaction
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		// 1- find the user that the token belongs token
		user, err := s.getUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}

		// 2- update the user
		user.IsActive = true
		if err := s.update(ctx, tx, user); err != nil {
			return err
		}

		// 3- clean the invitations
		if err := s.deleteUserInvitations(ctx, tx, user.ID); err != nil {
			return err
		}
		return nil
	})
}

func (s *UserStore) update(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `UPDATE users SET username = $1, email = $2, is_active = $3 WHERE id = $4`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	_, err := tx.ExecContext(ctx, query, user.Username, user.Email, user.IsActive, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, expiresIn time.Duration, userId int64) error {
	query := `INSERT INTO user_invitations (token, user_id, expiresIn) VALUES ($1,$2,$3)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, userId, time.Now().Add(expiresIn))
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.created_at, is_active
		FROM users u
		JOIN user_invitations ui ON u.id = ui.user_id
		WHERE ui.token = $1 AND ui.expiresin > $2;
	`

	// hashing the token
	hash := sha256.Sum256([]byte(token))
	hashedToken := hex.EncodeToString(hash[:])

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}
	err := tx.QueryRowContext(ctx, query, hashedToken, time.Now()).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.IsActive)
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

func (s *UserStore) deleteUserInvitations(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `DELETE FROM user_invitations WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) Delete(ctx context.Context, userId int64) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.delete(ctx, tx, userId); err != nil {
			return err
		}

		if err := s.deleteUserInvitations(ctx, tx, userId); err != nil {
			return err
		}
		return nil
	})
}

func (s *UserStore) delete(ctx context.Context, tx *sql.Tx, userId int64) error {
	query := `DELETE FROM users WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userId)
	if err != nil {
		return err
	}

	return nil
}
