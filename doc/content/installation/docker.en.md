---
title: Docker
description: Run yarr in a Docker container with a persistent volume.
weight: 3
---

Pull the image from the registry:

```sh
docker pull ghcr.io/nkanaev/yarr:latest
```

Run the container:

```sh
docker run -d --name yarr \
  -p 7070:7070 \
  -v yarr-data:/data \
  ghcr.io/nkanaev/yarr:latest -db /data/yarr.db
```

The server listens on `127.0.0.1:7070`. To expose it on all interfaces, add the `-addr` option:

```sh
docker run -d --name yarr \
  -p 7070:7070 \
  -v yarr-data:/data \
  ghcr.io/nkanaev/yarr:latest -db /data/yarr.db -addr 0.0.0.0:7070
```

The database is stored in the volume `yarr-data`.
