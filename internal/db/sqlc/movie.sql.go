// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: movie.sql

package db

import (
	"context"
)

const createMovie = `-- name: CreateMovie :one
INSERT INTO movies(title, director_id, rating, poster, summary)
VALUES($1, $2, $3, $4, $5)
RETURNING id, title, director_id, rating, poster, summary, created_at
`

type CreateMovieParams struct {
	Title      string `json:"title"`
	DirectorID int64  `json:"director_id"`
	Rating     int16  `json:"rating"`
	Poster     string `json:"poster"`
	Summary    string `json:"summary"`
}

func (q *Queries) CreateMovie(ctx context.Context, arg CreateMovieParams) (Movie, error) {
	row := q.db.QueryRowContext(ctx, createMovie,
		arg.Title,
		arg.DirectorID,
		arg.Rating,
		arg.Poster,
		arg.Summary,
	)
	var i Movie
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.DirectorID,
		&i.Rating,
		&i.Poster,
		&i.Summary,
		&i.CreatedAt,
	)
	return i, err
}

const deleteMovie = `-- name: DeleteMovie :exec
DELETE FROM movies
WHERE id = $1
`

func (q *Queries) DeleteMovie(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteMovie, id)
	return err
}

const getMovie = `-- name: GetMovie :one
SELECT id, title, director_id, rating, poster, summary, created_at
FROM movies
WHERE id = $1
ORDER BY id
LIMIT 1
`

func (q *Queries) GetMovie(ctx context.Context, id int64) (Movie, error) {
	row := q.db.QueryRowContext(ctx, getMovie, id)
	var i Movie
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.DirectorID,
		&i.Rating,
		&i.Poster,
		&i.Summary,
		&i.CreatedAt,
	)
	return i, err
}

const listMovies = `-- name: ListMovies :many
SELECT id, title, director_id, rating, poster, summary, created_at
FROM movies
ORDER BY id
LIMIT 8
`

func (q *Queries) ListMovies(ctx context.Context) ([]Movie, error) {
	rows, err := q.db.QueryContext(ctx, listMovies)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Movie{}
	for rows.Next() {
		var i Movie
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.DirectorID,
			&i.Rating,
			&i.Poster,
			&i.Summary,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}