-- name: GetChirps :many
SELECT
    chirps.id,
    chirps.created_at,
    chirps.updated_at,
    chirps.body,
    chirps.user_id
FROM
    chirps
ORDER BY 
    chirps.created_at ASC;