-- name: CreateURL :one
INSERT INTO urls (
    original_url,short_url,is_custom,expires_at
) VALUES (
        $1, $2, $3, $4
) 
RETURNING *;