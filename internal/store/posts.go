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
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Version   int       `json:"version"`
	Comments  []Comment `json:"comments"`
}

// instead of adding Comments in Post struct we can make a seperate struct like PostWithMetaData.

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
	SELECT id, user_id, title, content, created_at, updated_at, tags, version
	FROM posts
	WHERE id = $1
	`
	// We are not including comments(by using left join here) with this query because it might be a perfomance bottleneck. We might need post without comments in some cases.
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
		&post.Version,
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

func (s *PostStore) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM posts 
		WHERE id = $1
	`
	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostStore) Update(ctx context.Context, post *Post) error {
	query := `
		UPDATE posts
		SET title = $1 , content = $2, version = version + 1 
		WHERE id = $3 AND version = $4
		RETURNING version
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		post.ID,
		post.Version).Scan(&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}
