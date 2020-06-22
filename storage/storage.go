package storage

import (
	"os"
	"log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var initQuery string = `
create table if not exists folders (
 id             integer primary key autoincrement,
 title          text not null,
 is_expanded    boolean
);

create table if not exists feeds (
 id             integer primary key autoincrement,
 folder_id      references folders(id),
 title          text not null,
 description    text,
 link           text,
 feed_link      text,
 icon           text
);

create index if not exists idx_feed_folder_id on feeds(folder_id);

create table if not exists items (
 id             string primary key,
 feed_id        references feeds(id),
 title          text,
 link           text,
 description    text,
 content        text,
 author         text,
 date           integer,
 date_updated   integer,
 status         integer,
 image          text
);

create index if not exists idx_item_feed_id on items(feed_id);
create index if not exists idx_item_status  on items(status);
`

type Storage struct {
	db *sql.DB
	log *log.Logger
}

func New() (*Storage, error) {
	path := "./storage.db"
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(initQuery)
	if err != nil {
		return nil, err
	}
	logger := log.New(os.Stdout, "storage: ", log.Ldate | log.Ltime | log.Lshortfile)
	return &Storage{db: db, log: logger}, nil
}

func intOrNil(id int64) interface{} {
	if id == 0 {
		return nil
	}
	return id
}
