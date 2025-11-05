package models

import "time"

type RefreshToken struct {
	ID        uint
	Token     string
	UserId    string
	UserAgent string
	Ip        string
	CreatedAt time.Time
	ExpiresAt time.Time
}
