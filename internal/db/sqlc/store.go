package db

import "database/sql"

// Store provides all DB functions
type Store struct {
	*Queries
}

// NewStore creates a new store instance
func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
	}
}
