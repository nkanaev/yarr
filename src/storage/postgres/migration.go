package postgres

import (
	"database/sql"
	"fmt"
	"log"
)

var migrations = []func(*sql.Tx) error{
	m01_initial,
}

var maxVersion = int64(len(migrations))

func migrate(db *sql.DB) error {
	if _, err := db.Exec(
		`create table if not exists schema_version (version bigint primary key)`,
	); err != nil {
		return fmt.Errorf("create schema_version table: %w", err)
	}

	var version int64
	err := db.QueryRow(
		`select coalesce(max(version), 0) from schema_version`,
	).Scan(&version)
	if err != nil {
		return fmt.Errorf("read schema version: %w", err)
	}

	if version >= maxVersion {
		return nil
	}

	log.Printf("db version is %d. migrating to %d", version, maxVersion)

	for v := version + 1; v <= maxVersion; v++ {
		log.Printf("[migration:%d] starting", v)

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("migration %d begin tx: %w", v, err)
		}

		if err := migrations[v-1](tx); err != nil {
			tx.Rollback()
			return fmt.Errorf("migration %d: %w", v, err)
		}

		if _, err := tx.Exec(
			`insert into schema_version (version) values ($1)
			 on conflict do nothing`, v,
		); err != nil {
			tx.Rollback()
			return fmt.Errorf("migration %d record version: %w", v, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("migration %d commit: %w", v, err)
		}

		log.Printf("[migration:%d] done", v)
	}
	return nil
}

func m01_initial(tx *sql.Tx) error {
	_, err := tx.Exec(`
		create table if not exists folders (
			id          bigserial primary key,
			title       text not null,
			is_expanded boolean not null default false
		);
		create unique index if not exists idx_folder_title on folders(title);

		create table if not exists feeds (
			id          bigserial primary key,
			folder_id   bigint references folders(id) on delete set null,
			title       text not null,
			description text,
			link        text,
			feed_link   text not null,
			icon        bytea
		);
		create index if not exists idx_feed_folder_id on feeds(folder_id);
		create unique index if not exists idx_feed_feed_link on feeds(feed_link);

		create table if not exists items (
			id           bigserial primary key,
			guid         text not null,
			feed_id      bigint not null references feeds(id) on delete cascade,
			title        text,
			link         text,
			content      text,
			date         timestamptz,
			date_arrived timestamptz,
			last_arrived timestamptz,
			status       integer,
			media_links  jsonb,
			search       tsvector
		);
		create index if not exists idx_item_feed_id on items(feed_id);
		create index if not exists idx_item__date_id_status on items(date, id, status);
		create unique index if not exists idx_item_guid on items(feed_id, guid);
		create index if not exists idx_item_search on items using gin(search);

		create table if not exists settings (
			key text primary key,
			val jsonb
		);

		create table if not exists feed_states (
			feed_id        bigint primary key references feeds(id) on delete cascade,
			last_refreshed timestamptz not null default '1970-01-01 00:00:00+00',
			last_error     text not null default '',
			http_lmod      text not null default '',
			http_etag      text not null default ''
		);
	`)
	return err
}
