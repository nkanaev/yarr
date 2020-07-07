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
 is_expanded    boolean not null default false
);

create table if not exists feeds (
 id             integer primary key autoincrement,
 folder_id      references folders(id),
 title          text not null,
 description    text,
 link           text,
 feed_link      text not null,
 icon           text
);

create index if not exists idx_feed_folder_id on feeds(folder_id);
create unique index if not exists idx_feed_feed_link on feeds(feed_link);

create table if not exists items (
 id             integer primary key autoincrement,
 guid           string not null,
 feed_id        references feeds(id),
 title          text,
 link           text,
 description    text,
 content        text,
 author         text,
 date           datetime,
 date_updated   datetime,
 status         integer,
 image          text
);

create index if not exists idx_item_feed_id on items(feed_id);
create index if not exists idx_item_status  on items(status);
create unique index if not exists idx_item_guid on items(guid);

create table if not exists settings (
 key            string primary key,
 val            blob
);
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
	/*
	_, err = db.Exec(initQuery)
	if err != nil {
		return nil, err
	}
	*/
	logger := log.New(os.Stdout, "storage: ", log.Ldate | log.Ltime | log.Lshortfile)
	return &Storage{db: db, log: logger}, nil
}

func intOrNil(id int64) interface{} {
	if id == 0 {
		return nil
	}
	return id
}
