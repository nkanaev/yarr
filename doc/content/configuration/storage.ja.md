---
title: ストレージ
description: SQLite または PostgreSQL でデータを保存する方法。
weight: 4
---

yarr はすべてのデータをデータベースに保存します。`-db` フラグ（または `YARR_DB` 環境変数）を設定することで、SQLite と PostgreSQL のいずれかを選択できます。

## SQLite

SQLite がデフォルトであり、`-db` の値が通常のファイルパスである場合に使用されます。データベースパスが指定されない場合、yarr はユーザー設定ディレクトリ内の `storage.db` ファイルにデータを保存します:

- Linux - `$XDG_CONFIG_HOME/yarr/storage.db` (または `~/.config/yarr/storage.db`)
- macOS - `~/Library/Application Support/yarr/storage.db`
- Windows - `%AppData%\yarr\storage.db`

カスタム SQLite ファイルを使用する場合:

```sh
yarr -db /path/to/data.db
```

SQLite サポートは `mattn/go-sqlite3` ドライバーによって提供されています。ファイルパスには同ドライバーの[接続文字列引数](https://pkg.go.dev/github.com/mattn/go-sqlite3#readme-connection-string)を指定できます。引数が渡されない場合、yarr はデフォルト値 `_journal=WAL&_sync=NORMAL&_busy_timeout=5000&cache=shared` を適用します。独自の引数を渡すと、デフォルト値は完全に置き換えられます。

## PostgreSQL

`-db` の値が `postgres://` または `postgresql://` で始まる接続文字列の場合、PostgreSQL が有効になります。

```sh
yarr -db 'postgres://user:password@host:port/dbname?sslmode=disable'
```

データベースは事前に存在している必要があります。yarr は起動時にテーブルを自動的に作成・更新します。

PostgreSQL サポートは `lib/pq` ドライバーによって提供されています。接続文字列は URI パラメータをサポートしています。オプションの詳細については、[lib/pq 接続設定](https://pkg.go.dev/github.com/lib/pq#Config)を参照してください。