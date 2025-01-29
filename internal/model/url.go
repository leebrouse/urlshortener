package model

import "time"

type CreateURLRequest struct {
	Orginal_url string `json:"original_url" validate:"required,url"`
	CustomCode  string `json:"custom_code,omitempty" validate:"omitempty,alphanum,min=4,max=10"`
	Duration    *int   `json:"duration,omitempty" validate:"omitempty,min=1,max=100"`
}

type CreateUResponse struct {
	ShortUrl  string    `json:"short_url"`
	ExpiresAt time.Time `json:"expires_at"`
}
