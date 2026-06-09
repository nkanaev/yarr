package sqlite

import (
	"database/sql"
	"log"
	"strings"

	"github.com/mattn/go-sqlite3"
	"github.com/nkanaev/yarr/src/content/htmlutil"
)

func init() {
	sql.Register("sqlite3_yarr", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			return conn.RegisterFunc("strip_html", htmlutil.ExtractText, true)
		},
	})
}

type SQLiteStorage struct {
	db *sql.DB
}

func New(path string) (*SQLiteStorage, error) {
	if pos := strings.IndexRune(path, '?'); pos == -1 {
		params := "_journal=WAL&_sync=NORMAL&_busy_timeout=5000&cache=shared"
		log.Printf("opening db with params: %s", params)
		path = path + "?" + params
	}

	db, err := sql.Open("sqlite3_yarr", path)
	if err != nil {
		return nil, err
	}

	if err = migrate(db); err != nil {
		return nil, err
	}
	return &SQLiteStorage{db: db}, nil
}

func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}
