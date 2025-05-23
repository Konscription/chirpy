-- name: CreateUser :one
INSERT INTO users ( created_at, updated_at, email, hashed_password)
VALUES (
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING id, created_at, updated_at, email;

-- name: LookupUserbyEmail :one
SELECT
    id,
    created_at,
    updated_at,
    email,
    hashed_password
FROM
    users
WHERE
    email = $1;
