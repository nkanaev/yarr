---
title: Fever API
description: Connect Fever-compatible clients to yarr.
weight: 5
---

yarr supports the Fever API at the `/fever` endpoint. This lets
Fever-compatible third-party clients (e.g. Reeder, Unread) connect to your yarr
instance and read your feeds.

## Enable the Fever API

The Fever API requires authentication. Set a username and password on your yarr
server with the `-auth` flag (or the `YARR_AUTH` environment variable):

```sh
yarr -auth username:password
```

## Configure a Fever-compatible client

1. Make sure your yarr server is reachable from the client.
2. In the client, point it to your yarr server URL plus `/fever`.
3. Enter the username and password you configured on the yarr server.

## Notes

The Fever API specification is not precise, so server and client
implementations can have compatibility issues.

The following apps have been tested with yarr:

> Different apps accept different URL formats. Note whether the URL includes
> the `http://` scheme and a trailing `/`.

| App | Platforms | Config Server URL |
| :-- | :-- | :-- |
| [Reeder](https://reederapp.com/) | macOS, iOS | `127.0.0.1:7070/fever` or `http://127.0.0.1:7070/fever` |
| [ReadKit](https://readkit.app/) | macOS, iOS | `http://127.0.0.1:7070/fever` |
| [Fluent Reader](https://github.com/yang991178/fluent-reader) | macOS, Windows | `http://127.0.0.1:7070/fever/` |
| [Unread](https://www.goldenhillsoftware.com/unread/) | iOS | `http://127.0.0.1:7070/fever` |
| [Fiery Feeds](https://voidstern.net/fiery-feeds) | macOS, iOS | `http://127.0.0.1:7070/fever` |
