---
title: Docker
description: Запуск yarr в Docker-контейнере с сохранением данных в volume.
weight: 3
---

Скачайте образ из реестра:

```sh
docker pull ghcr.io/nkanaev/yarr:latest
```

Запустите контейнер:

```sh
docker run -d --name yarr \
  -p 7070:7070 \
  -v yarr-data:/data \
  ghcr.io/nkanaev/yarr:latest -db /data/yarr.db
```

Сервер слушает `127.0.0.1:7070`. Чтобы открыть доступ на всех интерфейсах, добавьте опцию `-addr`:

```sh
docker run -d --name yarr \
  -p 7070:7070 \
  -v yarr-data:/data \
  ghcr.io/nkanaev/yarr:latest -db /data/yarr.db -addr 0.0.0.0:7070
```

База данных хранится в томе `yarr-data`.