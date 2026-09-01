---
title: Docker
description: 使用持久卷在 Docker 容器中运行 yarr。
weight: 3
---

该镜像可从 [Docker Hub](https://hub.docker.com/r/nkanaev/yarr) 和 [GitHub Container Registry](https://github.com/nkanaev/yarr/pkgs/container/yarr) 获取。

从镜像仓库拉取镜像：

```sh
docker pull nkanaev/yarr:latest
```

运行容器：

```sh
docker run -d --name yarr \
  -p 7070:7070 \
  -v yarr-data:/data \
  nkanaev/yarr:latest -db /data/yarr.db -addr 0.0.0.0:7070
```

数据库存储在数据卷 `yarr-data` 中。