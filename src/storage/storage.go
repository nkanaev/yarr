package storage

import (
	"log"
	"database/sql"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

var sqlite3conn = []*sqlite3.SQLiteConn{}

func CopyDb(path string, from int, to int) (error) {

	db, err := sql.Open("sqlite3_with_con_hook", path)
	if err != nil {
		return err
	}
	db.Ping()

	bk, err := sqlite3conn[to].Backup("main", sqlite3conn[from], "main")
	if err != nil {
		return err
	}

	_, err = bk.Step(-1)
	if err != nil {
		return err
	}
	defer bk.Finish()
	defer db.Close()

	// drop last entry from cons
	sqlite3conn = sqlite3conn[:len(sqlite3conn) - 1]

	return err
}

func New(path string, db_fast bool) (*Storage, error) {
	var db * sql.DB
	var err error

	if db_fast {
		if len(sqlite3conn) == 0 {
			log.Printf("Registering sqlite3_with_con_hook ...")
			sql.Register("sqlite3_with_con_hook",
				&sqlite3.SQLiteDriver{
					ConnectHook: func(conn *sqlite3.SQLiteConn) error {
						// log.Printf("[sqlite3_with_con_hook] Appending conn ...")
						sqlite3conn = append(sqlite3conn, conn)
						return nil
					},
				})
		}
		db, err = sql.Open("sqlite3_with_con_hook", ":memory:")
	}else{
		db, err = sql.Open("sqlite3", path)
	}

	if err != nil {
		return nil, err
	}

	if db_fast {

		db.Ping()

		log.Printf("Loading %s to in-memory db", path)
		if err = CopyDb(path, 1, 0); err != nil {
			return nil, err
		}
	}

	// TODO: https://foxcpp.dev/articles/the-right-way-to-use-go-sqlite3
	db.SetMaxOpenConns(1)

	if err = migrate(db); err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func Close(store *Storage, path string, db_fast bool) error {
	var err error

	if db_fast {
		log.Printf("Saving in-memory db to %s", path)
		err = CopyDb(path, 0, 1)
	}
	defer store.db.Close()

	return err
}
