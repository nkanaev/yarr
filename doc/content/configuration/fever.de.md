---
title: Fever API
description: Verbinden Sie Fever-kompatible Clients mit yarr.
weight: 5
---

yarr unterstützt die Fever-API unter dem Endpunkt `/fever`. Dadurch können Fever-kompatible Drittanbieter-Clients (z. B. Reeder, Unread) eine Verbindung zu Ihrer yarr-Instanz herstellen und Ihre Feeds lesen.

## Aktivieren der Fever-API

Die Fever-API erfordert eine Authentifizierung. Legen Sie einen Benutzernamen und ein Passwort auf Ihrem yarr-Server mit dem Parameter `-auth` (oder der Umgebungsvariable `YARR_AUTH`) fest:

```sh
yarr -auth benutzername:passwort
```

## Konfigurieren eines Fever-kompatiblen Clients

1. Stellen Sie sicher, dass Ihr yarr-Server vom Client aus erreichbar ist.
2. Tragen Sie im Client die URL Ihres yarr-Servers gefolgt von `/fever` ein.
3. Geben Sie den Benutzernamen und das Passwort ein, die Sie auf dem yarr-Server konfiguriert haben.

## Hinweise

Die Spezifikation der Fever-API ist nicht exakt definiert, weshalb Kompatibilitätsprobleme zwischen Server und Client auftreten können.

Die folgenden Apps wurden mit yarr getestet:

> Verschiedene Apps akzeptieren unterschiedliche URL-Formate. Achten Sie darauf, ob die URL das Schema `http://` und einen schließenden Schrägstrich `/` enthält.

| App | Plattformen | Konfigurierte Server-URL |
| :-- | :-- | :-- |
| [Reeder](https://reederapp.com/) | macOS, iOS | `127.0.0.1:7070/fever` oder `http://127.0.0.1:7070/fever` |
| [ReadKit](https://readkit.app/) | macOS, iOS | `http://127.0.0.1:7070/fever` |
| [Fluent Reader](https://github.com/yang991178/fluent-reader) | macOS, Windows | `http://127.0.0.1:7070/fever/` |
| [Unread](https://www.goldenhillsoftware.com/unread/) | iOS | `http://127.0.0.1:7070/fever` |
| [Fiery Feeds](https://voidstern.net/fiery-feeds) | macOS, iOS | `http://127.0.0.1:7070/fever` |