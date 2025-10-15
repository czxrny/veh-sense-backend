package models

import "time"

type RefreshInfo struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	RefreshKey string    `json:"refresh_key"`
	ExpiresAt  time.Time `json:"expires_at"`
}

type TokenRefreshRequest struct {
	UserID     int    `json:"user_id" validate:"required"`
	RefreshKey string `json:"refresh_key" validate:"required"`
}
