---
title: Fever API
description: 将兼容 Fever 的客户端连接到 yarr。
weight: 5
---

yarr 在 `/fever` 端点支持 Fever API。这允许兼容 Fever 的第三方客户端（例如 Reeder、Unread）连接到您的 yarr 实例并读取您的 Feed。

## 启用 Fever API

Fever API 需要身份验证。在 yarr 服务器上使用 `-auth` 标志（或 `YARR_AUTH` 环境变量）设置用户名和密码：

{{< code "command-line-auth.sh" >}}

## 配置兼容 Fever 的客户端

1. 确保您的 yarr 服务器可从客户端访问。
2. 在客户端中，将服务器地址指向 yarr 服务器 URL 加 `/fever`。
3. 输入您在 yarr 服务器上配置的用户名和密码。

## 注意事项

Fever API 规范不够精确，因此服务器和客户端实现之间可能存在兼容性问题。

以下应用已通过 yarr 测试：

> 不同的应用接受不同的 URL 格式。请注意 URL 是否包含 `http://` 协议前缀以及结尾的 `/`。

| 应用 | 平台 | 配置服务器 URL |
| :-- | :-- | :-- |
| [Reeder](https://reederapp.com/) | macOS, iOS | `127.0.0.1:7070/fever` 或 `http://127.0.0.1:7070/fever` |
| [ReadKit](https://readkit.app/) | macOS, iOS | `http://127.0.0.1:7070/fever` |
| [Fluent Reader](https://github.com/yang991178/fluent-reader) | macOS, Windows | `http://127.0.0.1:7070/fever/` |
| [Unread](https://www.goldenhillsoftware.com/unread/) | iOS | `http://127.0.0.1:7070/fever` |
| [Fiery Feeds](https://voidstern.net/fiery-feeds) | macOS, iOS | `http://127.0.0.1:7070/fever` |