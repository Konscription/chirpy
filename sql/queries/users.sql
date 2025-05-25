-- name: CreateUser :one
INSERT INTO users ( created_at, updated_at, email, hashed_password)
VALUES (
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING id, created_at, updated_at, email, is_chirpy_red;

-- name: LookupUserbyEmail :one
SELECT
    id,
    created_at,
    updated_at,
    email,
    hashed_password,
    is_chirpy_red
FROM
    users
WHERE
    email = $1;

-- name: UpdateUser :exec
UPDATE users
SET
    updated_at = NOW(),
    hashed_password = $2,
    email = $3
WHERE
    id = $1;

-- name: LookupUserById :one
SELECT
    id,
    created_at,
    updated_at,
    email,
    is_chirpy_red
FROM
    users
WHERE
    id = $1;