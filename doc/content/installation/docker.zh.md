---
title: Docker
description: 使用持久卷在 Docker 容器中运行 yarr。
weight: 3
---

从镜像仓库拉取镜像：

```sh
docker pull ghcr.io/nkanaev/yarr:latest
```

运行容器：

```sh
docker run -d --name yarr \
  -p 7070:7070 \
  -v yarr-data:/data \
  ghcr.io/nkanaev/yarr:latest -db /data/yarr.db
```

服务器监听 `127.0.0.1:7070`。要暴露在所有接口上，请添加 `-addr` 选项：

```sh
docker run -d --name yarr \
  -p 7070:7070 \
  -v yarr-data:/data \
  ghcr.io/nkanaev/yarr:latest -db /data/yarr.db -addr 0.0.0.0:7070
```

数据库存储在数据卷 `yarr-data` 中。