package repository

import (
	models "Server_part_finance_control/server/internal/domains"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

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
		INSERT INTO users (id, email, username, password_hash, vk_id, created_at)
        VALUES ($1, LOWER($2), $3, $4, $5, $6)
        RETURNING id, created_at	
	`

	err := r.db.QueryRow(ctx, query, user.Email, user.Username, user.VKID, user.CreatedAt, password).Scan(&user.ID, &user.CreatedAt)
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

func (r *UserRepository) GetUserByVKID(ctx context.Context, vkID string) (*models.User, error) {
	const query = `
		SELECT id, email, username, password_hash, created_at 
		FROM users 
		WHERE vk_id = $1
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, vkID).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, fmt.Errorf("user with VK ID %s not found", vkID)
	case err != nil:
		return nil, fmt.Errorf("failed to get user by VK ID: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) UpdateUserVKID(ctx context.Context, userID, vkID string) error {
	const query = `
		UPDATE users 
		SET vk_id = $1 
		WHERE id = $2
	`

	_, err := r.db.Exec(ctx, query, vkID, userID)
	if err != nil {
		return fmt.Errorf("failed to update VK ID: %w", err)
	}

	return nil
}

// Ping проверяет соединение с базой данных
func (r *UserRepository) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return r.db.Ping(ctx)
}
