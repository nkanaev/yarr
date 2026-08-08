---
title: Docker
description: Führen Sie yarr in einem Docker-Container mit einem dauerhaften Volume aus.
weight: 3
---

Laden Sie das Image aus der Registry herunter:

```sh
docker pull ghcr.io/nkanaev/yarr:latest
```

Starten Sie den Container:

```sh
docker run -d --name yarr \
  -p 7070:7070 \
  -v yarr-data:/data \
  ghcr.io/nkanaev/yarr:latest -db /data/yarr.db
```

Der Server lauscht auf `127.0.0.1:7070`. Um ihn auf allen Schnittstellen erreichbar zu machen, fügen Sie die Option `-addr` hinzu:

```sh
docker run -d --name yarr \
  -p 7070:7070 \
  -v yarr-data:/data \
  ghcr.io/nkanaev/yarr:latest -db /data/yarr.db -addr 0.0.0.0:7070
```

Die Datenbank wird im Volume `yarr-data` gespeichert.