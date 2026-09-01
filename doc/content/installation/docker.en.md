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

## Docker Compose with PostgreSQL

To run yarr with PostgreSQL, use a `docker-compose.yml` file:

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

Start the services:

```sh
docker compose up -d
```
