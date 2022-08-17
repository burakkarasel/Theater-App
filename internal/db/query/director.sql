-- name: CreateDirector :one
INSERT INTO directors(first_name, last_name, oscars)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetDirector :one
SELECT *
FROM directors
WHERE id = $1
LIMIT 1;

-- name: ListDirectors :many
SELECT *
FROM directors
ORDER BY id
LIMIT $1
OFFSET $2;