package model

import "time"

type CreateURLRequest struct {
	OriginalURL string `json:"original_url" validate:"required,url"`
	CustomCode  string `json:"custom_code,omitempty" validate:"omitempty,alphanum,min=4,max=10"`
	Duration    *int   `json:"duration,omitempty" validate:"omitempty,min=1,max=100"`
}

type CreateURLResponse struct {
	ShortUrl  string    `json:"short_url"`
	ExpiresAt time.Time `json:"expires_at"`
}
