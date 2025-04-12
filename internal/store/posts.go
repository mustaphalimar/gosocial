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
	Version   int       `json:"version"`
}

type FeedPost struct {
	Post
	CommentsCount int `json:"comments_count"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
	INSERT INTO posts(content,title,user_id,tags)
	VALUES ($1,$2,$3,$4) RETURNING id,created_at,updated_at`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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
	query := `SELECT id,user_id,title,content,created_at,updated_at,tags,version FROM posts WHERE id=$1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var post Post
	err := s.db.QueryRowContext(ctx, query, postId).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		pq.Array(&post.Tags),
		&post.Version,
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

func (s *PostStore) Update(ctx context.Context, post *Post) error {
	query := `
		UPDATE posts SET title = $1,content = $2,version = version + 1 WHERE id = $3 AND version = $4 RETURNING version;
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, post.Title, post.Content, post.ID, post.Version).Scan(&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrorNotFound
		default:
			return err
		}
	}

	return nil
}

func (s *PostStore) Delete(ctx context.Context, postId int64) error {
	query := `
		DELETE FROM posts where id = $1;
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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

func (s *PostStore) GetUserFeed(ctx context.Context, userId int64, fq PaginatedFeedQuery) ([]FeedPost, error) {
	query := `
    select p.id,p.user_id,p.title,p.content,p.created_at,p.version,p.tags,u.username, count(c.id) as comments_count
    from posts p
    LEFT JOIN comments c ON c.post_id = p.id
    LEFT JOIN users u ON p.user_id = u.id
    JOIN followers f ON f.follower_id = p.user_id or p.user_id = $1
    where f.user_id = $1 or p.user_id = $1
    group by p.id, u.id
    order by p.created_at ` + fq.Sort + `
    limit $2 offset $3
    `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userId, fq.Limit, fq.Offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var feedPosts []FeedPost
	for rows.Next() {
		var p FeedPost
		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
			&p.Version,
			pq.Array(&p.Tags),
			&p.User.Username,
			&p.CommentsCount,
		)

		if err != nil {
			return nil, err
		}
		feedPosts = append(feedPosts, p)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return feedPosts, nil
}
