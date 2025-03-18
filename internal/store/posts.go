package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Comments  []Comment `json:"comments"`
	User      User      `json:"user"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
	INSERT INTO posts(content,title,user_id,tags)
	VALUES ($1,$2,$3,$4) RETURNING id,created_at,updated_at`

	err := s.db.QueryRowContext(ctx, query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (s *PostStore) GetById(ctx context.Context, postId int64) (*Post, error) {
	query := `SELECT id,user_id,title,content,created_at,updated_at, tags FROM posts WHERE id=$1`

	var post Post
	err := s.db.QueryRowContext(ctx, query, postId).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		pq.Array(&post.Tags),
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrorNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}

// func (s *PostStore) Update(ctx context.Context, postId int64)  (*Post,error) {
// 	query := `
// 		UPDATE posts SET title = 'My Very 1st post' WHERE posts.id = 1;
// 	`

// 	var post *Post
// 	err := s.db.QueryRowContext(ctx, query,
// 		post.Content,
// 		post.Title,
// 		post.UserID,
// 		pq.Array(post.Tags),
// 	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

// 	if err != nil {
// 		return nil,err
// 	}
// 	return post,err

// }

func (s *PostStore) Delete(ctx context.Context, postId int64) error {
	query := `
		DELETE FROM posts where id = $1;
	`

	res, err := s.db.ExecContext(ctx, query, postId)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()

	if rows == 0 {
		return ErrorNotFound
	}

	return nil
}
