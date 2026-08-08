---
title: Stockage
description: Comment stocker vos données avec SQLite ou PostgreSQL.
weight: 4
---

yarr stocke toutes les données dans une base de données. Vous pouvez choisir entre SQLite et PostgreSQL en définissant l'option `-db` (ou la variable d'environnement `YARR_DB`).

## SQLite

SQLite est utilisé par défaut lorsque la valeur de `-db` est un chemin de fichier ordinaire. Si aucun chemin n'est fourni, yarr stocke les données dans un fichier `storage.db` situé dans le dossier de configuration utilisateur :

- Linux - `$XDG_CONFIG_HOME/yarr/storage.db` (ou `~/.config/yarr/storage.db`)
- macOS - `~/Library/Application Support/yarr/storage.db`
- Windows - `%AppData%\yarr\storage.db`

Pour utiliser un fichier SQLite personnalisé :

```sh
yarr -db /chemin/vers/donnees.db
```

La prise en charge de SQLite repose sur le pilote `mattn/go-sqlite3`. Le chemin du fichier accepte les [arguments de chaîne de connexion](https://pkg.go.dev/github.com/mattn/go-sqlite3#readme-connection-string) de ce pilote. Si aucun argument n'est fourni, yarr applique les valeurs par défaut `_journal=WAL&_sync=NORMAL&_busy_timeout=5000&cache=shared`. Si vous passez vos propres arguments, ils remplacent entièrement les paramètres par défaut.

## PostgreSQL

PostgreSQL est activé lorsque la valeur de `-db` est une chaîne de connexion commençant par `postgres://` ou `postgresql://`.

```sh
yarr -db 'postgres://utilisateur:motdepasse@hote:port/nombd?sslmode=disable'
```

La base de données doit déjà exister ; yarr crée et met à jour les tables automatiquement au démarrage.

La prise en charge de PostgreSQL repose sur le pilote `lib/pq`. La chaîne de connexion prend en charge les paramètres URI. Pour une liste complète des options, consultez la [configuration de connexion lib/pq](https://pkg.go.dev/github.com/lib/pq#Config).