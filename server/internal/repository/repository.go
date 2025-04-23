package repository

import (
	models "Server_part_finance_control/server/internal/domains"
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	db *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) *UserRepository{
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User, password string) error{
	query := `
		INSERT INTO users (email, username, password_hash)
		VALUES (LOWER($1), $2, $3)
		RETURNING id, created_at	
	`

	err := r.db.QueryRow(ctx, query, user.Email, user.Username, password).Scan(&user.ID, &user.CreatedAt)
	if err  != nil {
		if strings.Contains(err.Error(), "duplicated key value") {
			return fmt.Errorf("user with this email already exists")
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	query := `
		SELECT id, email, username, password_hash, created_at
		FROM users
		WHERE email = LOWER($1)
	`

	err := r.db.QueryRow(ctx, query, strings.ToLower(email)).Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows{
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}