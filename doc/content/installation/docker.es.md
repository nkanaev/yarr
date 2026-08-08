---
title: Docker
description: Ejecuta yarr en un contenedor Docker con un volumen persistente.
weight: 3
---

Obtén la imagen desde el registro:

```sh
docker pull ghcr.io/nkanaev/yarr:latest
```

Ejecuta el contenedor:

```sh
docker run -d --name yarr \
  -p 7070:7070 \
  -v yarr-data:/data \
  ghcr.io/nkanaev/yarr:latest -db /data/yarr.db
```

El servidor escucha en `127.0.0.1:7070`. Para exponerlo en todas las interfaces, añade la opción `-addr`:

```sh
docker run -d --name yarr \
  -p 7070:7070 \
  -v yarr-data:/data \
  ghcr.io/nkanaev/yarr:latest -db /data/yarr.db -addr 0.0.0.0:7070
```

La base de datos se almacena en el volumen `yarr-data`.