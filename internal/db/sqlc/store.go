package db

import "database/sql"

// Store provides all DB functions
type Store interface {
	Querier
}

// Store provides all DB functions
type SQLStore struct {
	*Queries
}

// NewStore creates a new store instance
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db),
	}
}
