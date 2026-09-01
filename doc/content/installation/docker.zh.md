---
title: Docker
description: 使用持久卷在 Docker 容器中运行 yarr。
weight: 3
---

该镜像可从 [Docker Hub](https://hub.docker.com/r/nkanaev/yarr) 和 [GitHub Container Registry](https://github.com/nkanaev/yarr/pkgs/container/yarr) 获取。

从镜像仓库拉取镜像：

{{< code "docker-pull.sh" >}}

运行容器：

{{< code "docker-run.sh" >}}

数据库存储在数据卷 `yarr-data` 中。

## Docker Compose 与 PostgreSQL

要结合 PostgreSQL 运行 yarr，请使用 `docker-compose.yml` 文件：

{{< code "docker-compose.yml" >}}

启动服务：

{{< code "docker-compose-up.sh" >}}
