---
title: Speicher
description: So speichern Sie Ihre Daten mit SQLite oder PostgreSQL.
weight: 4
---

yarr speichert alle Daten in einer Datenbank. Sie können zwischen SQLite und PostgreSQL wählen, indem Sie den Parameter `-db` (oder die Umgebungsvariable `YARR_DB`) setzen.

## SQLite

SQLite ist der Standard und wird verwendet, wenn der Wert von `-db` ein gewöhnlicher Dateipfad ist. Wenn kein Datenbankpfad angegeben ist, speichert yarr die Daten in einer Datei `storage.db` im Konfigurationsverzeichnis des Benutzers:

- Linux - `$XDG_CONFIG_HOME/yarr/storage.db` (oder `~/.config/yarr/storage.db`)
- macOS - `~/Library/Application Support/yarr/storage.db`
- Windows - `%AppData%\yarr\storage.db`

So verwenden Sie eine benutzerdefinierte SQLite-Datei:

```sh
yarr -db /pfad/zu/daten.db
```

Die SQLite-Unterstützung basiert auf dem Treiber `mattn/go-sqlite3`. Der Dateipfad akzeptiert die [Verbindungszeichenfolgen-Argumente](https://pkg.go.dev/github.com/mattn/go-sqlite3#readme-connection-string) dieses Treibers. Wenn keine Argumente angegeben werden, wendet yarr die Standardwerte `_journal=WAL&_sync=NORMAL&_busy_timeout=5000&cache=shared` an. Wenn Sie eigene Argumente übergeben, ersetzen diese die Standardwerte vollständig.

## PostgreSQL

PostgreSQL wird aktiviert, wenn der Wert von `-db` eine Verbindungszeichenfolge ist, die mit `postgres://` oder `postgresql://` beginnt.

```sh
yarr -db 'postgres://benutzer:passwort@host:port/datenbankname?sslmode=disable'
```

Die Datenbank muss bereits existieren; yarr erstellt und aktualisiert die Tabellen beim Start automatisch.

Die PostgreSQL-Unterstützung basiert auf dem Treiber `lib/pq`. Die Verbindungszeichenfolge unterstützt URI-Parameter. Eine vollständige Liste der Optionen finden Sie in der [lib/pq Verbindungs-Konfiguration](https://pkg.go.dev/github.com/lib/pq#Config).