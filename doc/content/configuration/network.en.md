---
title: Network
description: Where the server listens — TCP addresses and Unix sockets.
weight: 6
---

Use the `-addr` flag (or the `YARR_ADDR` environment variable) to choose where the server listens.

By default, the server listens on `127.0.0.1:7070`, which means only the machine it runs on can reach it.
Pass `0.0.0.0` as the host to listen on all interfaces for external access:

{{< code "command-line-addr-wildcard.sh" >}}

You can also listen on a Unix socket instead of TCP. Prefix the socket path with `unix:`, which is handy when the server sits behind a web server (e.g. nginx or Caddy) that proxies to it:

{{< code "command-line-addr-socket.sh" >}}
