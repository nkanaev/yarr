---
title: 选项
description: 命令行标志及其对应的环境变量。
weight: 3
---

服务器接受命令行参数和/或环境变量作为配置选项。命令行标志的优先级高于对应的环境变量。

| 标志 | 环境变量 | 描述 |
| ------------ | -------------------- | ---------------------------------------------------------------------------- |
| `-addr` | `YARR_ADDR` | 服务器运行地址（默认 `127.0.0.1:7070`） |
| `-base` | `YARR_BASE` | 服务 URL 的基础路径 |
| `-auth` | `YARR_AUTH` | 格式为 `username:password` 的用户名和密码 |
| `-auth-file` | `YARR_AUTHFILE` | 包含 `username:password` 的文件路径。优先级高于 `-auth` |
| `-cert-file` | `YARR_CERTFILE` | TLS 证书文件路径 |
| `-key-file` | `YARR_KEYFILE` | TLS 密钥文件路径 |
| `-db` | `YARR_DB` | 存储数据库文件路径 |
| `-log-file` | `YARR_LOGFILE` | 日志文件路径 |
| `-open` | — | 在浏览器中打开服务器 |

## HTTPS

同时需要 `-cert-file` 和 `-key-file` 才能启用 HTTPS。