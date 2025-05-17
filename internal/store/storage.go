package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetById(context.Context, int64) (*Post, error)
		Update(context.Context, *Post) error
		Delete(context.Context, int64) error
		GetUserFeed(context.Context, int64, PaginatedFeedQuery) ([]FeedPost, error)
		DeleteAll(context.Context) error
	}
	Users interface {
		GetById(context.Context, int64) (*User, error)
		GetByEmail(context.Context, string) (*User, error)
		Create(context.Context, *sql.Tx, *User) error
		CreateAndInvite(ctx context.Context, user *User, token string, expiresIn time.Duration) error
		DeleteAll(context.Context) error
		Activate(ctx context.Context, token string) error
		Delete(ctx context.Context, userId int64) error
	}
	Comments interface {
		Create(context.Context, *Comment) error
		GetByPostId(context.Context, int64) ([]Comment, error)
		DeleteAll(context.Context) error
	}
	Followers interface {
		Follow(ctx context.Context, userToFollow int64, followingUser int64) error
		Unfollow(ctx context.Context, userToUnfollow int64, unfollowingUser int64) error
	}
}

var (
	ErrorNotFound        = errors.New("Record not found.")
	QueryTimeoutDuration = time.Second * 5
	ErrConflict          = errors.New("Resource already exists")
	ErrDuplicateEmail    = errors.New("Email already in use")
	ErrDuplicateUsername = errors.New("Username already in use")
)

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostStore{db},
		Users:     &UserStore{db},
		Comments:  &CommentStore{db},
		Followers: &FollowerStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
