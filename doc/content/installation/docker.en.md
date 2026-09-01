---
title: Docker
description: Run yarr in a Docker container with a persistent volume.
weight: 3
---

The image is available on [Docker Hub](https://hub.docker.com/r/nkanaev/yarr) and [GitHub Container Registry](https://github.com/nkanaev/yarr/pkgs/container/yarr).

Pull the image from the registry:

```sh
docker pull nkanaev/yarr:latest
```

Run the container:

```sh
docker run -d --name yarr \
  -p 7070:7070 \
  -v yarr-data:/data \
  nkanaev/yarr:latest -db /data/yarr.db -addr 0.0.0.0:7070
```

The database is stored in the volume `yarr-data`.
