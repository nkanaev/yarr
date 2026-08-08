---
title: Armazenamento
description: Como armazenar os seus dados com SQLite ou PostgreSQL.
weight: 4
---

O yarr armazena todos os dados numa base de dados. Pode escolher entre SQLite e PostgreSQL definindo a opção `-db` (ou a variável de ambiente `YARR_DB`).

## SQLite

O SQLite é a opção predefinida e é utilizado quando o valor de `-db` é um caminho de ficheiro ordinário. Se não for especificado nenhum caminho, o yarr armazena os dados num ficheiro `storage.db` dentro da pasta de configuração do utilizador:

- Linux - `$XDG_CONFIG_HOME/yarr/storage.db` (ou `~/.config/yarr/storage.db`)
- macOS - `~/Library/Application Support/yarr/storage.db`
- Windows - `%AppData%\yarr\storage.db`

Para utilizar um ficheiro SQLite personalizado:

```sh
yarr -db /caminho/para/dados.db
```

O suporte para SQLite é fornecido pelo controlador `mattn/go-sqlite3`. O caminho do ficheiro aceita os [argumentos de string de ligação](https://pkg.go.dev/github.com/mattn/go-sqlite3#readme-connection-string) desse controlador. Quando não são fornecidos argumentos, o yarr aplica as predefinições `_journal=WAL&_sync=NORMAL&_busy_timeout=5000&cache=shared`. Se passar os seus próprios argumentos, estes substituem totalmente as predefinições.

## PostgreSQL

O PostgreSQL é ativado quando o valor de `-db` é uma string de ligação que começa por `postgres://` ou `postgresql://`.

```sh
yarr -db 'postgres://utilizador:palavrapasse@anfitriao:porta/nomedabd?sslmode=disable'
```

A base de dados já deve existir; o yarr cria e atualiza as tabelas automaticamente ao iniciar.

O suporte para PostgreSQL é fornecido pelo controlador `lib/pq`. A string de ligação suporta parâmetros URI. Para consultar a lista completa de opções, consulte a [configuração de ligação do lib/pq](https://pkg.go.dev/github.com/lib/pq#Config).