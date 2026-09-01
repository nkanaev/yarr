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

## Docker Compose 与 PostgreSQL

要结合 PostgreSQL 运行 yarr，请使用 `docker-compose.yml` 文件：

```yaml
services:
  yarr:
    image: nkanaev/yarr:latest
    command: ["-db", "postgres://yarr:yarr@postgres:5432/yarr?sslmode=disable", "-addr", "0.0.0.0:7070"]
    ports:
      - "7070:7070"
    depends_on:
      postgres:
        condition: service_healthy

  postgres:
    image: postgres:17-alpine
    environment:
      POSTGRES_USER: yarr
      POSTGRES_PASSWORD: yarr
      POSTGRES_DB: yarr
    volumes:
      - postgres-data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U yarr -d yarr"]
      interval: 5s
      timeout: 5s
      retries: 5

volumes:
  postgres-data:
```

启动服务：

```sh
docker compose up -d
```
