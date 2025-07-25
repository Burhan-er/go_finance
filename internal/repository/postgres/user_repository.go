package postgres

import (
	"context"
	"database/sql"
	"go_finance/internal/domain"
	"time"
)

// --- User Repository ---

type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository.
func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db: db}
}

// CreateUser inserts a new user record into the database.
func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (username, email, password_hash, role, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	// Set a default role if not specified.
	if user.Role == "" {
		user.Role = domain.UserRole
	}

	err := r.db.QueryRowContext(ctx, query, user.Username, user.Email, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
	return err
}

// GetUserByEmail retrieves a user by their email address.
func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, role, created_at, updated_at FROM users WHERE email = $1`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Or a custom "not found" error
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByID retrieves a user by their unique ID.
func (r *userRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	// Corrected SQL query (removed extra comma)
	query := `SELECT id, username, email, password_hash, role, created_at, updated_at FROM users WHERE id = $1`

	var user domain.User
	// Correctly use Scan to get the row data
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Or a custom "not found" error
		}
		return nil, err
	}

	return &user, nil
}

// GetAllUsers retrieves all users from the database.
func (r *userRepository) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	// Corrected SQL query (removed extra comma and WHERE clause)
	query := `SELECT id, username, email, password_hash, role, created_at, updated_at FROM users`

	// Use QueryContext for multiple rows
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	// Loop through all the returned rows
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


	// Check for errors during row iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

