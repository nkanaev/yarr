---
title: 命令行
description: 将 yarr 作为服务器运行，并从任何浏览器访问。
weight: 2
---

## 下载

从[发布页面](https://github.com/nkanaev/yarr/releases/latest)下载命令行版本。文件名不包含 `_gui`。

| 操作系统 | 架构 | 下载 |
|---|---|---|
| macOS | Apple Silicon | [yarr_darwin_arm64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_darwin_arm64.zip) |
| macOS | Intel | [yarr_darwin_amd64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_darwin_amd64.zip) |
| Windows | ARM64 | [yarr_windows_arm64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_windows_arm64.zip) |
| Windows | x86-64 | [yarr_windows_amd64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_windows_amd64.zip) |
| Linux | ARM64 | [yarr_linux_arm64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_linux_arm64.zip) |
| Linux | ARMv7 | [yarr_linux_armv7.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_linux_armv7.zip) |
| Linux | x86-64 | [yarr_linux_amd64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_linux_amd64.zip) |

## 运行

启动服务器：

```sh
./yarr
```

服务器默认监听 `127.0.0.1:7070`。yarr 会自动将数据存储在用户配置文件夹中。

## 示例

在所有网络接口上运行服务器并启用密码保护：

```sh
yarr -addr 0.0.0.0:7070 -auth alice:secret
```

在浏览器中打开 `http://host:7070`，使用 `alice` / `secret` 登录。