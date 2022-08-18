-- name: CreateTicket :one
INSERT INTO tickets(movie_id, ticket_owner, child, adult, total)
VALUES($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetTicket :one
SELECT *
FROM tickets
WHERE id = $1
LIMIT 1;

-- name: ListTickets :many
SELECT *
FROM tickets
WHERE ticket_owner = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: DeleteTickets :exec
DELETE FROM tickets
WHERE id = $1;