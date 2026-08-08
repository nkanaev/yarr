---
title: コマンドライン
description: yarr をサーバーとして実行し、任意のブラウザからアクセスします。
weight: 2
---

## ダウンロード

[リリースページ](https://github.com/nkanaev/yarr/releases/latest)からコマンドライン版をダウンロードします。ファイル名には `_gui` が含まれていません。

| OS | アーキテクチャ | ダウンロード |
|---|---|---|
| macOS | Apple Silicon | [yarr_darwin_arm64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_darwin_arm64.zip) |
| macOS | Intel | [yarr_darwin_amd64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_darwin_amd64.zip) |
| Windows | ARM64 | [yarr_windows_arm64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_windows_arm64.zip) |
| Windows | x86-64 | [yarr_windows_amd64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_windows_amd64.zip) |
| Linux | ARM64 | [yarr_linux_arm64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_linux_arm64.zip) |
| Linux | ARMv7 | [yarr_linux_armv7.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_linux_armv7.zip) |
| Linux | x86-64 | [yarr_linux_amd64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_linux_amd64.zip) |

## 実行

サーバーを起動します:

```sh
./yarr
```

デフォルトではサーバーは `127.0.0.1:7070` で待機します。yarr はユーザー設定フォルダにデータを自動的に保存します。

## 実行例

パスワード保護を有効にして、すべてのネットワークインターフェースでサーバーを実行する例:

```sh
yarr -addr 0.0.0.0:7070 -auth alice:secret
```

ブラウザで `http://host:7070` を開き、ユーザー名 `alice` / パスワード `secret` でログインします。