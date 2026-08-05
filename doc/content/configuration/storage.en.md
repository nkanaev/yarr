---
title: Storage
description: How to store your data with SQLite or PostgreSQL.
weight: 4
---

yarr stores all data in a database. You can choose between SQLite and
PostgreSQL by setting the `-db` flag (or the `YARR_DB` environment variable).

## SQLite

SQLite is the default and is used when the `-db` value is an ordinary file path.
If no database path is given, yarr stores the data in a `storage.db` file inside
the user's config directory:

- Linux - `$XDG_CONFIG_HOME/yarr/storage.db` (or `~/.config/yarr/storage.db`)
- macOS - `~/Library/Application Support/yarr/storage.db`
- Windows - `%AppData%\yarr\storage.db`

To use a custom SQLite file:

```sh
yarr -db /path/to/data.db
```

SQLite support is backed by the `mattn/go-sqlite3` driver. The file path accepts
the [connection string arguments](https://pkg.go.dev/github.com/mattn/go-sqlite3#readme-connection-string)
of that driver. When no arguments are given, yarr applies the defaults
`_journal=WAL&_sync=NORMAL&_busy_timeout=5000&cache=shared`. If you pass your
own arguments, they replace the defaults entirely.

## PostgreSQL

PostgreSQL is enabled when the `-db` value is a connection string that starts
with `postgres://` or `postgresql://`.

```sh
yarr -db 'postgres://user:password@host:port/dbname?sslmode=disable'
```

The database must already exist; yarr creates and updates the tables
automatically on startup.

PostgreSQL support is backed by the `lib/pq` driver. The connection string
supports URI parameters. For a full list of options, see the
[lib/pq connection configuration](https://pkg.go.dev/github.com/lib/pq#Config).
