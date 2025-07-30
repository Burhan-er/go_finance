package postgres

import (
	"context"
	"database/sql"
	"errors"
	"go_finance/internal/domain"
	"time"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (username, email, password_hash, role, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if user.Role == "" {
		user.Role = domain.UserRole
	}

	err := r.db.QueryRowContext(ctx, query, user.Username, user.Email, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
	return err
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, role, created_at, updated_at FROM users WHERE email = $1`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil 
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, role, created_at, updated_at FROM users WHERE id = $1`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	query := `SELECT id, username, email, password_hash, role, created_at, updated_at FROM users`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		if err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.PasswordHash,
			&user.Role, &user.CreatedAt, &user.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) UpdateUserByID(ctx context.Context, id string) (*domain.User, error) {
	query := `UPDATE users SET username = COALESCE(NULLIF($2, ''), username), email = COALESCE(NULLIF($3, ''), email), updated_at = $4 WHERE id = $1`

	user, _ := r.GetUserByID(ctx, id)
	_, err := r.db.ExecContext(ctx, query, id, user.Username, user.Email, time.Now())
	if err != nil {
		return nil, err
	}
	return user, err
}

func (r *userRepository) DeleteUserByID(ctx context.Context, id string) (*domain.User, error) {
	query := `DELETE FROM users WHERE id = $1`
	user, _ := r.GetUserByID(ctx, id)
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	return user, err
}
