-- name: ListMovies :many
SELECT *
FROM movies
ORDER BY id
LIMIT 8;

-- name: GetMovie :one
SELECT *
FROM movies
WHERE id = $1
ORDER BY id
LIMIT 1;

-- name: CreateMovie :one
INSERT INTO movies(title, director_id, rating, poster)
VALUES($1, $2, $3, $4)
RETURNING *;

-- name: DeleteMovie :exec
DELETE FROM movies
WHERE id = $1;