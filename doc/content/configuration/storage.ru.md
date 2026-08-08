---
title: Хранилище
description: Как хранить данные с помощью SQLite или PostgreSQL.
weight: 4
---

yarr хранит все данные в базе данных. Вы можете выбрать между SQLite и
PostgreSQL, указав флаг `-db` (или переменную окружения `YARR_DB`).

## SQLite

SQLite используется по умолчанию, если значение `-db` является обычным путём к файлу.
Если путь к базе данных не указан, yarr сохраняет данные в файле `storage.db` внутри
папки конфигурации пользователя:

- Linux — `$XDG_CONFIG_HOME/yarr/storage.db` (или `~/.config/yarr/storage.db`)
- macOS — `~/Library/Application Support/yarr/storage.db`
- Windows — `%AppData%\yarr\storage.db`

Чтобы использовать собственный файл SQLite:

```sh
yarr -db /path/to/data.db
```

Поддержка SQLite реализована на драйвере `mattn/go-sqlite3`. Путь к файлу принимает
[аргументы строки подключения](https://pkg.go.dev/github.com/mattn/go-sqlite3#readme-connection-string)
этого драйвера. Если аргументы не переданы, yarr применяет значения по умолчанию
`_journal=WAL&_sync=NORMAL&_busy_timeout=5000&cache=shared`. Если вы передаёте
свои аргументы, они полностью заменяют значения по умолчанию.

## PostgreSQL

PostgreSQL используется, если значение `-db` является строкой подключения, начинающейся
с `postgres://` или `postgresql://`.

```sh
yarr -db 'postgres://user:password@host:port/dbname?sslmode=disable'
```

База данных уже должна существовать; yarr создаёт и обновляет таблицы
автоматически при запуске.

Поддержка PostgreSQL реализована на драйвере `lib/pq`. Строка подключения
поддерживает параметры URI. Полный список опций см. в разделе
[конфигурации подключения lib/pq](https://pkg.go.dev/github.com/lib/pq#Config).