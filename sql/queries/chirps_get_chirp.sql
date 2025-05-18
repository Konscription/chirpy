-- name: GetChirp :one
SELECT
    chirps.id,
    chirps.created_at,
    chirps.updated_at,
    chirps.body,
    chirps.user_id
FROM
    chirps
WHERE
    chirps.id = $1;
