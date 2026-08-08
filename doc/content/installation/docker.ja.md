---
title: Docker
description: 永続ボリュームを使用して Docker コンテナ内で yarr を実行します。
weight: 3
---

レジストリからイメージを取得します:

```sh
docker pull ghcr.io/nkanaev/yarr:latest
```

コンテナを実行します:

```sh
docker run -d --name yarr \
  -p 7070:7070 \
  -v yarr-data:/data \
  ghcr.io/nkanaev/yarr:latest -db /data/yarr.db
```

サーバーは `127.0.0.1:7070` で待機します。すべてのインターフェースに公開するには、`-addr` オプションを追加します:

```sh
docker run -d --name yarr \
  -p 7070:7070 \
  -v yarr-data:/data \
  ghcr.io/nkanaev/yarr:latest -db /data/yarr.db -addr 0.0.0.0:7070
```

データベースはボリューム `yarr-data` に保存されます。