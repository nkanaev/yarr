---
title: Almacenamiento
description: Cómo almacenar tus datos con SQLite o PostgreSQL.
weight: 4
---

yarr almacena todos los datos en una base de datos. Puedes elegir entre SQLite y PostgreSQL configurando el parámetro `-db` (o la variable de entorno `YARR_DB`).

## SQLite

SQLite es la opción por defecto y se utiliza cuando el valor de `-db` es una ruta de archivo convencional. Si no se especifica ninguna ruta, yarr guarda los datos en un archivo `storage.db` dentro de la carpeta de configuración del usuario:

- Linux - `$XDG_CONFIG_HOME/yarr/storage.db` (o `~/.config/yarr/storage.db`)
- macOS - `~/Library/Application Support/yarr/storage.db`
- Windows - `%AppData%\yarr\storage.db`

Para usar un archivo SQLite personalizado:

```sh
yarr -db /ruta/a/datos.db
```

El soporte para SQLite está respaldado por el controlador `mattn/go-sqlite3`. La ruta de archivo acepta los [argumentos de la cadena de conexión](https://pkg.go.dev/github.com/mattn/go-sqlite3#readme-connection-string) de dicho controlador. Cuando no se proporcionan argumentos, yarr aplica por defecto `_journal=WAL&_sync=NORMAL&_busy_timeout=5000&cache=shared`. Si pasas tus propios argumentos, estos reemplazarán completamente los valores por defecto.

## PostgreSQL

PostgreSQL se habilita cuando el valor de `-db` es una cadena de conexión que comienza por `postgres://` o `postgresql://`.

```sh
yarr -db 'postgres://usuario:contraseña@host:puerto/nombrebd?sslmode=disable'
```

La base de datos ya debe existir; yarr crea y actualiza las tablas automáticamente al iniciar.

El soporte para PostgreSQL está respaldado por el controlador `lib/pq`. La cadena de conexión admite parámetros URI. Para ver la lista completa de opciones, consulta la [configuración de conexión de lib/pq](https://pkg.go.dev/github.com/lib/pq#Config).