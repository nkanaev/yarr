---
title: Docker
description: Execute o yarr num contentor Docker com um volume persistente.
weight: 3
---

Obtenha a imagem a partir do registo:

```sh
docker pull ghcr.io/nkanaev/yarr:latest
```

Execute o contentor:

```sh
docker run -d --name yarr \
  -p 7070:7070 \
  -v yarr-data:/data \
  ghcr.io/nkanaev/yarr:latest -db /data/yarr.db
```

O servidor escuta em `127.0.0.1:7070`. Para o expor em todas as interfaces, adicione a opção `-addr`:

```sh
docker run -d --name yarr \
  -p 7070:7070 \
  -v yarr-data:/data \
  ghcr.io/nkanaev/yarr:latest -db /data/yarr.db -addr 0.0.0.0:7070
```

A base de dados é armazenada no volume `yarr-data`.