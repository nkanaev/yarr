package storage

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	// TODO: https://foxcpp.dev/articles/the-right-way-to-use-go-sqlite3
	db.SetMaxOpenConns(1)

	if err = migrate(db); err != nil {
		return nil, err
	}
	return &Storage{db: db}, nil
}
