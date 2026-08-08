---
title: Fever API
description: Fever 互換クライアントを yarr に接続します。
weight: 5
---

yarr は `/fever` エンドポイントで Fever API をサポートしています。これにより、Fever 互換のサードパーティ製クライアント（例: Reeder、Unread）から yarr インスタンスに接続してフィードを閲覧できます。

## Fever API の有効化

Fever API には認証が必要です。yarr サーバーで `-auth` フラグ（または `YARR_AUTH` 環境変数）を使用してユーザー名とパスワードを設定します:

```sh
yarr -auth username:password
```

## Fever 互換クライアントの設定

1. クライアントから yarr サーバーにアクセスできることを確認します。
2. クライアントで、yarr サーバーの URL に `/fever` を加えたアドレスを指定します。
3. yarr サーバーで設定したユーザー名とパスワードを入力します。

## 備考

Fever API の仕様は厳密ではないため、サーバーとクライアントの実装間において互換性の問題が生じる場合があります。

以下のアプリは yarr で動作確認されています:

> アプリによって受け入れられる URL 形式が異なります。URL に `http://` スキームや末尾の `/` が含まれているかに注意してください。

| アプリ | プラットフォーム | 設定サーバー URL |
| :-- | :-- | :-- |
| [Reeder](https://reederapp.com/) | macOS, iOS | `127.0.0.1:7070/fever` または `http://127.0.0.1:7070/fever` |
| [ReadKit](https://readkit.app/) | macOS, iOS | `http://127.0.0.1:7070/fever` |
| [Fluent Reader](https://github.com/yang991178/fluent-reader) | macOS, Windows | `http://127.0.0.1:7070/fever/` |
| [Unread](https://www.goldenhillsoftware.com/unread/) | iOS | `http://127.0.0.1:7070/fever` |
| [Fiery Feeds](https://voidstern.net/fiery-feeds) | macOS, iOS | `http://127.0.0.1:7070/fever` |