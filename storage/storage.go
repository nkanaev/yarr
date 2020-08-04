package storage

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var initQuery string = `
create table if not exists folders (
 id             integer primary key autoincrement,
 title          text not null,
 is_expanded    boolean not null default false
);

create unique index if not exists idx_folder_title on folders(title);

create table if not exists feeds (
 id             integer primary key autoincrement,
 folder_id      references folders(id),
 title          text not null,
 description    text,
 link           text,
 feed_link      text not null,
 icon           blob
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
 image          text,
 search_rowid   integer
);

create index if not exists idx_item_feed_id on items(feed_id);
create index if not exists idx_item_status  on items(status);
create unique index if not exists idx_item_guid on items(guid);
create index if not exists idx_item_search_rowid on items(search_rowid);

create table if not exists settings (
 key            string primary key,
 val            blob
);

create virtual table if not exists search using fts4(title, description, content);

create trigger if not exists del_item_search after delete on items begin
  delete from search where rowid = old.search_rowid;
end;
`

type Storage struct {
	db  *sql.DB
	log *log.Logger
}

func New(path string, logger *log.Logger) (*Storage, error) {
	initialize := false
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			initialize = true
		} else {
			return nil, err
		}
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if initialize {
		if _, err := db.Exec(initQuery); err != nil {
			return nil, err
		}
	}
	return &Storage{db: db, log: logger}, nil
}

func intOrNil(id int64) interface{} {
	if id == 0 {
		return nil
	}
	return id
}
