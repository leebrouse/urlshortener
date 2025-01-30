-- name: CreateURL :one
INSERT INTO urls (
    original_url,short_url,is_custom,expires_at
) VALUES (
        $1, $2, $3, $4
) 
RETURNING *;

-- name: IsShortCodeAvailable :one
SELECT NOT EXISTS (
    SELECT 1 FROM urls 
    WHERE short_url = $1
) AS available;

-- name: GetURLByShortCode :one
SELECT * FROM urls 
WHERE short_url = $1
AND expires_at > current_timestamp;

-- name: DeleteURLExpired :exec
DELETE FROM urls 
WHERE expires_at <= current_timestamp;

