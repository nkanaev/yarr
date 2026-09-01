---
title: Docker
description: Run yarr in a Docker container with a persistent volume.
weight: 3
---

The image is available on [Docker Hub](https://hub.docker.com/r/nkanaev/yarr) and [GitHub Container Registry](https://github.com/nkanaev/yarr/pkgs/container/yarr).

Pull the image from the registry:

{{< code "docker-pull.sh" >}}

Run the container:

{{< code "docker-run.sh" >}}

The database is stored in the volume `yarr-data`.

## Docker Compose with PostgreSQL

To run yarr with PostgreSQL, use a `docker-compose.yml` file:

{{< code "docker-compose.yml" >}}

Start the services:

{{< code "docker-compose-up.sh" >}}
