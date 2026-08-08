---
title: Optionen
description: Befehlszeilenparameter und ihre entsprechenden Umgebungsvariablen.
weight: 3
---

Der Server akzeptiert Optionen als Befehlszeilenargumente und/oder Umgebungsvariablen. Ein Befehlszeilenparameter hat Vorrang vor der entsprechenden Umgebungsvariable.

| Parameter | Umgebungsvariable | Beschreibung |
| ------------ | -------------------- | ---------------------------------------------------------------------------- |
| `-addr` | `YARR_ADDR` | Adresse, auf der der Server läuft (Standard `127.0.0.1:7070`) |
| `-base` | `YARR_BASE` | Basispfad der Dienst-URL |
| `-auth` | `YARR_AUTH` | Benutzername und Passwort im Format `benutzername:passwort` |
| `-auth-file` | `YARR_AUTHFILE` | Pfad zu einer Datei mit `benutzername:passwort`. Hat Vorrang vor `-auth` |
| `-cert-file` | `YARR_CERTFILE` | Pfad zur TLS-Zertifikatsdatei |
| `-key-file` | `YARR_KEYFILE` | Pfad zur TLS-Schlüsseldatei |
| `-db` | `YARR_DB` | Pfad zur Speicherdatei |
| `-log-file` | `YARR_LOGFILE` | Pfad zur Protokolldatei (Logfile) |
| `-open` | — | Server im Browser öffnen |

## HTTPS

Sowohl `-cert-file` als auch `-key-file` sind erforderlich, um HTTPS zu aktivieren.