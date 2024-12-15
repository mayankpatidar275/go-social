package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

// This is like model(in mvc). can be kept in a seperate model folder
// Here I kept the model tight with the storage fetching

type Post struct {
	ID        int64    `json:"id"`
	Content   string   `json:"content"`
	Title     string   `json:"title"`
	UserID    int64    `json:"user_id"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type PostStore struct {
	db *sql.DB
}

// Note: We can use ORM (more friendly) like GORM to avoid writing sql
// sqlx, sqlboiler are other libraries to make life easier

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (content, title, user_id, tags)
		VALUES ($1, $2, $3, $4) RETURNING Id, created_at, updated_at
	`
	// Note: Placeholders like $1 ensure:
	// The database driver treats the inputs as data only, not as part of the SQL query.
	// Malicious inputs can’t "break out" of the query and execute harmful commands.

	// pq.Array(post.Tags) converts the Go []string (slice) into a PostgreSQL array, which is the expected format for the tags column
	// $1 corresponds to the first argument after ctx and query
	err := s.db.QueryRowContext(
		ctx, query, post.Content, post.Title, post.UserID, pq.Array(post.Tags),
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
	// post.ID is an integer value.
	// &post.ID is the address of the ID field within the Post struct.
	// Note: &post will be the address of the pointer post
	// fmt.Println(post.ID)   // Equivalent to (*post).ID // pointer to pointer rarely used

	if err != nil {
		return err
	}

	return nil
}

// I tried it by sending the post pointer as argument instead of returning the pointer by creating the post here.
func (s *PostStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	// instead of this SELECT * FROM posts mention everything explicitly is better instead of implicit.
	query := `
	SELECT id, user_id, title, content, created_at, updated_at, tags
	FROM posts
	WHERE id = $1
	`
	var post Post
	if err := s.db.QueryRowContext(ctx, query, id).Scan(
		//should be in same order as select query
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		pq.Array(&post.Tags),
	); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &post, nil
}
