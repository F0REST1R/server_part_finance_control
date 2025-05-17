package models

import "time"

type User struct {
	ID           string
	Email        string
	Username     string
	CreatedAt    time.Time
	VKID *string
	PasswordHash string
}