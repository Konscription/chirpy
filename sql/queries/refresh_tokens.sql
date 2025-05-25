-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3
)
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT
    u.id,
    u.created_at,
    u.updated_at,
    u.email,
    u.hashed_password,
    u.is_chirpy_red
FROM users u 
JOIN refresh_tokens rt ON u.id = rt.user_id
WHERE rt.token = $1
AND revoked_at IS NULL
AND rt.expires_at > NOW();

-- name: RevokeRefreshToken :one
UPDATE refresh_tokens
SET revoked_at = NOW(),
    updated_at = NOW()
WHERE token = $1
RETURNING *;