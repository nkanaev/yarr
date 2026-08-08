---
title: オプション
description: コマンドラインフラグとそれに対応する環境変数。
weight: 3
---

サーバーはコマンドライン引数および環境変数としてオプションを受け入れます。コマンドラインフラグは対応する環境変数よりも優先されます。

| フラグ | 環境変数 | 説明 |
| ------------ | -------------------- | ---------------------------------------------------------------------------- |
| `-addr` | `YARR_ADDR` | サーバーを実行するアドレス（デフォルト: `127.0.0.1:7070`） |
| `-base` | `YARR_BASE` | サービス URL のベースパス |
| `-auth` | `YARR_AUTH` | `username:password` 形式のユーザー名とパスワード |
| `-auth-file` | `YARR_AUTHFILE` | `username:password` が含まれるファイルへのパス（`-auth` より優先） |
| `-cert-file` | `YARR_CERTFILE` | TLS 証明書ファイルへのパス |
| `-key-file` | `YARR_KEYFILE` | TLS 秘密鍵ファイルへのパス |
| `-db` | `YARR_DB` | ストレージファイルへのパス |
| `-log-file` | `YARR_LOGFILE` | ログファイルへのパス |
| `-open` | — | ブラウザでサーバーを開く |

## HTTPS

HTTPS を有効にするには、`-cert-file` と `-key-file` の両方が必要です。