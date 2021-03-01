package storage

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type Storage struct {
	db  *sql.DB
	log *log.Logger
}

func New(path string, log *log.Logger) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)

	if err = migrate(db, log); err != nil {
		return nil, err
	}
	return &Storage{db: db, log: log}, nil
}
