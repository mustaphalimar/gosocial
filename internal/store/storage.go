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
	}
	Users interface {
		GetById(context.Context, int64) (*User, error)
		Create(context.Context, *User) error
	}
	Comments interface {
		Create(context.Context, *Comment) error
		GetByPostId(context.Context, int64) ([]Comment, error)
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
)

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostStore{db},
		Users:     &UserStore{db},
		Comments:  &CommentStore{db},
		Followers: &FollowerStore{db},
	}
}
