package storage

import (
	"database/sql"
	"fmt"
	"log"
)

var migrations = []func(*sql.Tx) error{
	m01_initial,
	m02_feed_states_and_errors,
	m03_on_delete_actions,
	m04_item_podcasturl,
	m05_move_description_to_content,
	m06_fill_missing_dates,
	m07_add_feed_size,
}

var maxVersion = int64(len(migrations))

func migrate(db *sql.DB) error {
	var version int64
	db.QueryRow("pragma user_version").Scan(&version)

	if version >= maxVersion {
		return nil
	}

	log.Printf("db version is %d. migrating to %d", version, maxVersion)

	for v := version + 1; v <= maxVersion; v++ {
		// Migrations altering schema using a sequence of steps due to SQLite limitations.
		// Must come with `pragma foreign_key_check` at the end. See:
		// "Making Other Kinds Of Table Schema Changes"
		// https://www.sqlite.org/lang_altertable.html
		trickyAlteration := (v == 3)

		log.Printf("[migration:%d] starting", v)

		if trickyAlteration {
			db.Exec("pragma foreign_keys=off;")
		}

		err := migrateVersion(v, db)

		if trickyAlteration {
			db.Exec("pragma foreign_keys=on;")
		}

		if err != nil {
			return err
		}

		log.Printf("[migration:%d] done", v)
	}
	return nil
}

func migrateVersion(v int64, db *sql.DB) error {
	var err error
	var tx *sql.Tx
	migratefunc := migrations[v-1]
	if tx, err = db.Begin(); err != nil {
		log.Printf("[migration:%d] failed to start transaction", v)
		return err
	}
	if err = migratefunc(tx); err != nil {
		log.Printf("[migration:%d] failed to migrate", v)
		tx.Rollback()
		return err
	}
	if _, err = tx.Exec(fmt.Sprintf("pragma user_version = %d", v)); err != nil {
		log.Printf("[migration:%d] failed to bump version", v)
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		log.Printf("[migration:%d] failed to commit changes", v)
		return err
	}
	return nil
}

func m01_initial(tx *sql.Tx) error {
	sql := `
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
		 date_arrived   datetime,
		 status         integer,
		 image          text,
		 search_rowid   integer
		);

		create index if not exists idx_item_feed_id on items(feed_id);
		create index if not exists idx_item_status  on items(status);
		create index if not exists idx_item_search_rowid on items(search_rowid);
		create unique index if not exists idx_item_guid on items(feed_id, guid);

		create table if not exists settings (
		 key            string primary key,
		 val            blob
		);

		create virtual table if not exists search using fts4(title, description, content);

		create trigger if not exists del_item_search after delete on items begin
		  delete from search where rowid = old.search_rowid;
		end;
	`
	_, err := tx.Exec(sql)
	return err
}

func m02_feed_states_and_errors(tx *sql.Tx) error {
	sql := `
		create table if not exists http_states (
		 feed_id        references feeds(id) unique,
		 last_refreshed datetime not null,

		 -- http header fields --
		 last_modified  string not null,
		 etag           string not null
		);

		create table if not exists feed_errors (
		 feed_id        references feeds(id) unique,
		 error          string
		);
	`
	_, err := tx.Exec(sql)
	return err
}

func m03_on_delete_actions(tx *sql.Tx) error {
	sql := `
		-- 01. create altered tables
		create table if not exists new_feeds (
		 id             integer primary key autoincrement,
		 folder_id      references folders(id) on delete set null,
		 title          text not null,
		 description    text,
		 link           text,
		 feed_link      text not null,
		 icon           blob
		);
		create table if not exists new_items (
		 id             integer primary key autoincrement,
		 guid           string not null,
		 feed_id        references feeds(id) on delete cascade,
		 title          text,
		 link           text,
		 description    text,
		 content        text,
		 author         text,
		 date           datetime,
		 date_updated   datetime,
		 date_arrived   datetime,
		 status         integer,
		 image          text,
		 search_rowid   integer
		);
		create table if not exists new_http_states (
		 feed_id        references feeds(id) on delete cascade unique,
		 last_refreshed datetime not null,
		 last_modified  string not null,
		 etag           string not null
		);
		create table if not exists new_feed_errors (
		 feed_id        references feeds(id) on delete cascade unique,
		 error          string
		);

		-- 02. transfer data into new tables
		insert into new_feeds select * from feeds;
		insert into new_items select * from items;
		insert into new_http_states select * from http_states;
		insert into new_feed_errors select * from feed_errors;

		-- 03. drop old tables
		drop table feeds;
		drop table items;
		drop table http_states;
		drop table feed_errors;

		-- 04. rename new tables
		alter table new_feeds rename to feeds;
		alter table new_items rename to items;
		alter table new_http_states rename to http_states;
		alter table new_feed_errors rename to feed_errors;

		-- 05. reconstruct indexes & triggers
		create index if not exists idx_feed_folder_id on feeds(folder_id);
		create unique index if not exists idx_feed_feed_link on feeds(feed_link);
		create index if not exists idx_item_feed_id on items(feed_id);
		create index if not exists idx_item_status  on items(status);
		create index if not exists idx_item_search_rowid on items(search_rowid);
		create unique index if not exists idx_item_guid on items(feed_id, guid);
		create trigger if not exists del_item_search after delete on items begin
		  delete from search where rowid = old.search_rowid;
		end;

		-- 06. check consistency
		pragma foreign_key_check;
	`
	_, err := tx.Exec(sql)
	return err
}

func m04_item_podcasturl(tx *sql.Tx) error {
	sql := `
		alter table items add column podcast_url text
	`
	_, err := tx.Exec(sql)
	return err
}

func m05_move_description_to_content(tx *sql.Tx) error {
	sql := `
		update items set content=description
		where length(content) = 0 or content is null
	`
	_, err := tx.Exec(sql)
	return err
}

func m06_fill_missing_dates(tx *sql.Tx) error {
	sql := `
		update items set date = 0 where date is null;
	`
	_, err := tx.Exec(sql)
	return err
}

func m07_add_feed_size(tx *sql.Tx) error {
	sql := `
		create table if not exists feed_sizes (
		 feed_id        references feeds(id) on delete cascade unique,
		 size           integer not null default 0
		);
	`
	_, err := tx.Exec(sql)
	return err
}
