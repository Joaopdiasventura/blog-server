package post

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostRepository interface {
	Create(ctx context.Context, createPostDTO CreatePostDTO) (Post, error)
	FindById(ctx context.Context, id string) (Post, error)
	FindMany(ctx context.Context, limit, offset int) ([]Post, error)
	FindAllByAuthor(ctx context.Context, author string) ([]Post, error)
	Delete(ctx context.Context, id string) error
}

type PostgresPostRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresPostRepository(pool *pgxpool.Pool) *PostgresPostRepository {
	return &PostgresPostRepository{pool: pool}
}

func (r *PostgresPostRepository) Create(ctx context.Context, createPostDTO CreatePostDTO) (Post, error) {
	query := `
		INSERT INTO posts (title, content, author_id)
		VALUES ($1, $2, $3)
		RETURNING id, title, content, author_id, created_at
	`
	var post Post

	err := r.pool.QueryRow(ctx, query, createPostDTO.Title, createPostDTO.Content, createPostDTO.AuthorId).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.AuthorId,
		&post.CreatedAt,
	)

	if err != nil {
		return Post{}, err
	}

	return post, nil
}

func (r *PostgresPostRepository) FindById(ctx context.Context, id string) (Post, error) {
	query := `
		SELECT id, title, content, author_id, created_at
		FROM posts
		WHERE id = $1
	`
	var post Post

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.AuthorId,
		&post.CreatedAt,
	)

	if err != nil {
		return Post{}, err
	}

	return post, nil
}

func (r *PostgresPostRepository) FindMany(ctx context.Context, limit, offset int) ([]Post, error) {
	query := `
		SELECT 
			p.id, 
			p.title, 
			p.created_at,
			u.id,
			u.name,
			u.email 
		FROM posts p
		JOIN users u ON p.author_id = u.id
		ORDER BY p.created_at DESC
		LIMIT $1 OFFSET $2
	`

	var posts []Post

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.CreatedAt,
			&post.Author.Id,
			&post.Author.Name,
			&post.Author.Email,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostgresPostRepository) FindAllByAuthor(ctx context.Context, author string) ([]Post, error) {
	query := `
		SELECT id, title, created_at
		FROM posts
		WHERE author_id = $1
		ORDER BY created_at DESC
	`
	var posts []Post

	rows, err := r.pool.Query(ctx, query, author)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostgresPostRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM posts
		WHERE id = $1
	`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}
