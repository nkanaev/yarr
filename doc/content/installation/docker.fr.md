---
title: Docker
description: Exécuter yarr dans un conteneur Docker avec un volume persistant.
weight: 3
---

Récupérez l'image depuis le registre :

```sh
docker pull ghcr.io/nkanaev/yarr:latest
```

Exécutez le conteneur :

```sh
docker run -d --name yarr \
  -p 7070:7070 \
  -v yarr-data:/data \
  ghcr.io/nkanaev/yarr:latest -db /data/yarr.db
```

Le serveur écoute sur `127.0.0.1:7070`. Pour l'exposer sur toutes les interfaces, ajoutez l'option `-addr` :

```sh
docker run -d --name yarr \
  -p 7070:7070 \
  -v yarr-data:/data \
  ghcr.io/nkanaev/yarr:latest -db /data/yarr.db -addr 0.0.0.0:7070
```

La base de données est stockée dans le volume `yarr-data`.