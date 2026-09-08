---
title: 网络
description: 服务器监听地址 —— TCP 地址与 Unix socket。
weight: 6
---

使用 `-addr` 标志（或 `YARR_ADDR` 环境变量）来选择服务器监听的地址。

默认情况下，服务器监听 `127.0.0.1:7070`，这意味着只有运行它的机器才能访问。
将主机设置为 `0.0.0.0` 可以监听所有网络接口以供外部访问：

{{< code "command-line-addr-wildcard.sh" >}}

您也可以监听 Unix socket 而不是 TCP。在 socket 路径前加上 `unix:` 前缀，这在服务器位于为其反向代理的 Web 服务器（例如 nginx 或 Caddy）之后时非常方便：

{{< code "command-line-addr-socket.sh" >}}
