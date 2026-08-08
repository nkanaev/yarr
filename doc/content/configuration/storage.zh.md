---
title: 存储
description: 如何使用 SQLite 或 PostgreSQL 存储数据。
weight: 4
---

yarr 将所有数据存储在数据库中。您可以通过设置 `-db` 标志（或 `YARR_DB` 环境变量）来选择 SQLite 或 PostgreSQL。

## SQLite

SQLite 是默认选项，当 `-db` 的值为普通文件路径时使用。如果没有提供数据库路径，yarr 会将数据存储在用户配置目录下的 `storage.db` 文件中：

- Linux - `$XDG_CONFIG_HOME/yarr/storage.db`（或 `~/.config/yarr/storage.db`）
- macOS - `~/Library/Application Support/yarr/storage.db`
- Windows - `%AppData%\yarr\storage.db`

要使用自定义 SQLite 文件：

```sh
yarr -db /path/to/data.db
```

SQLite 支持由 `mattn/go-sqlite3` 驱动提供。文件路径支持该驱动的[连接字符串参数](https://pkg.go.dev/github.com/mattn/go-sqlite3#readme-connection-string)。未提供参数时，yarr 默认使用 `_journal=WAL&_sync=NORMAL&_busy_timeout=5000&cache=shared`。如果您传递自定义参数，将完全覆盖默认参数。

## PostgreSQL

当 `-db` 的值是以 `postgres://` 或 `postgresql://` 开头的连接字符串时，启用 PostgreSQL。

```sh
yarr -db 'postgres://user:password@host:port/dbname?sslmode=disable'
```

数据库必须事先存在；yarr 会在启动时自动创建并更新表结构。

PostgreSQL 支持由 `lib/pq` 驱动提供。连接字符串支持 URI 参数。有关选项的完整列表，请参阅 [lib/pq 连接配置](https://pkg.go.dev/github.com/lib/pq#Config)。