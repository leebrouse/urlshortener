// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package repo

import (
	"context"
	"time"
)

const createURL = `-- name: CreateURL :one
INSERT INTO urls (
    original_url,short_url,is_custom,expires_at
) VALUES (
        $1, $2, $3, $4
) 
RETURNING id, original_url, short_url, is_custom, expires_at, created_at
`

type CreateURLParams struct {
	OriginalUrl string    `json:"original_url"`
	ShortUrl    string    `json:"short_url"`
	IsCustom    bool      `json:"is_custom"`
	ExpiresAt   time.Time `json:"expires_at"`
}

func (q *Queries) CreateURL(ctx context.Context, arg CreateURLParams) (Url, error) {
	row := q.db.QueryRowContext(ctx, createURL,
		arg.OriginalUrl,
		arg.ShortUrl,
		arg.IsCustom,
		arg.ExpiresAt,
	)
	var i Url
	err := row.Scan(
		&i.ID,
		&i.OriginalUrl,
		&i.ShortUrl,
		&i.IsCustom,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}
