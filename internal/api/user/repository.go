package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, createUserDTO CreateUserDTO) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	FindById(ctx context.Context, id string) (User, error)
	Update(ctx context.Context, id string, updateUserDTO UpdateUserDTO) (User, error)
	Delete(ctx context.Context, id string) error
}

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(poll *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{pool: poll}
}

func (r *PostgresUserRepository) Create(ctx context.Context, createUserDTO CreateUserDTO) (User, error) {
	query := `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, name, email, password
	`
	var user User

	err := r.pool.QueryRow(ctx, query, createUserDTO.Name, createUserDTO.Email, createUserDTO.Password).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (User, error) {
	query := `
		SELECT id, name, email, password
		FROM users
		WHERE email = $1
	`

	var user User

	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (r *PostgresUserRepository) FindById(ctx context.Context, id string) (User, error) {
	query := `
		SELECT id, name, email
		FROM users
		WHERE id = $1
	`

	var user User

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
	)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, id string, updateUserDTO UpdateUserDTO) (User, error) {
	query := `
		UPDATE users
		SET name = COALESCE($2, name),
		    email = COALESCE($3, email),
		    password = COALESCE($4, password)
		WHERE id = $1
		RETURNING id, name, email, password
	`

	var user User

	err := r.pool.QueryRow(ctx, query, id, updateUserDTO.Name, updateUserDTO.Email, updateUserDTO.Password).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}