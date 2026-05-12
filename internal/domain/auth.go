package domain

import "time"

type UserID string

type User struct {
	ID           UserID
	Email        string
	PasswordHash string
	Name         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}
