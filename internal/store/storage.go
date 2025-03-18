package store

import (
	"context"
	"database/sql"
	"errors"
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetById(context.Context, int64) (*Post, error)
		// Update(context.Context, int64,) (*Post,error)
		Delete(context.Context, int64) error
	}
	Users interface {
		Create(context.Context, *User) error
	}

	Comments interface {
		GetByPostId(context.Context, int64) ([]Comment, error)
	}
}

var (
	ErrorNotFound = errors.New("Record not found.")
)

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{db},
		Users:    &UserStore{db},
		Comments: &CommentStore{db},
	}
}
