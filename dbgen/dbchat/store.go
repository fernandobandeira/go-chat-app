package dbchat

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Store interface {
	Querier
	WithTx(*sql.Tx) Store
}

type SQLStore struct {
	db *sql.DB
	*Queries
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (s *SQLStore) WithTx(tx *sql.Tx) Store {
	return &SQLStore{
		db:      s.db,
		Queries: s.Queries.WithTx(tx),
	}
}
