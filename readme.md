# yarr

**yarr** (yet another rss reader) is a web-based feed aggregator which can be used both
as a desktop application and a personal self-hosted server.

The app is a single binary with an embedded database (SQLite).

![screenshot](etc/promo.png)

## usage

The latest prebuilt binaries for Linux/MacOS/Windows are available
[here](https://github.com/nkanaev/yarr/releases/latest).
The archives follow the naming convention `yarr_{OS}_{ARCH}[_gui].zip`, where:

* `OS` is the target operating system
* `ARCH` is the CPU architecture (`arm64` for AArch64, `amd64` for X86-64)
* `-gui` indicates that the binary ships with the GUI (tray icon), and is a command line application if omitted

Usage instructions:

* MacOS: place `yarr.app` in `/Applications` folder, [open the app][macos-open], click the anchor menu bar icon, select "Open".

* Windows: open `yarr.exe`, click the anchor system tray icon, select "Open".

* Linux: place `yarr` in `$HOME/.local/bin` and run [the script](etc/install-linux.sh).

[macos-open]: https://support.apple.com/en-gb/guide/mac-help/mh40616/mac

For self-hosting, see `yarr -h` for auth, tls & server configuration flags.

## deploying with once

yarr is compatible with [Basecamp Once](https://github.com/basecamp/once). To deploy with authentication:

```sh
docker run -d \
  -p 80:80 \
  -v yarr-data:/storage \
  -e SECRET_KEY_BASE="your-secret-key" \
  -e YARR_AUTH="username:password" \
  ghcr.io/sroberts/yarr:once-latest
```

Once mounts a persistent volume at `/storage` automatically. The database is stored at `/storage/db/yarr.db`.

### environment variables

| Variable | Description |
|---|---|
| `YARR_AUTH` | Username and password in `username:password` format |
| `SECRET_KEY_BASE` | Secret key for session token signing (provided by Once) |
| `DISABLE_SSL` | Set to `true` when running behind a reverse proxy that handles TLS (set by default in the Once image) |

The Once image listens on port 80 and includes backup/restore hooks for safe SQLite snapshots. A health check is available at `GET /up`.

### pre-built images

Pre-built multi-arch images (amd64/arm64) are published to GitHub Container Registry:

```sh
# Once-compatible image (port 80, /storage volume)
docker pull ghcr.io/sroberts/yarr:once-latest

# Standard image (port 7070, /data volume)
docker pull ghcr.io/sroberts/yarr:latest
```

### building the docker image locally

To build the Once-compatible image from source:

```sh
docker build -f etc/dockerfile.once -t yarr:once .
```

To build the standard image:

```sh
docker build -f etc/dockerfile -t yarr .
```

Run the locally built Once image:

```sh
docker run -d \
  -p 8080:80 \
  -v yarr-data:/storage \
  -e YARR_AUTH="username:password" \
  -e DISABLE_SSL="true" \
  yarr:once
```

### pushing to ghcr

To build and push your own image to GitHub Container Registry:

```sh
# Log in to GHCR
echo $GITHUB_TOKEN | docker login ghcr.io -u YOUR_USERNAME --password-stdin

# Build and tag for GHCR
docker build -f etc/dockerfile.once -t ghcr.io/YOUR_USERNAME/yarr:once-latest .

# Push
docker push ghcr.io/YOUR_USERNAME/yarr:once-latest
```

For multi-arch builds (amd64 + arm64), use Docker Buildx:

```sh
docker buildx create --use
docker buildx build -f etc/dockerfile.once \
  --platform linux/amd64,linux/arm64 \
  -t ghcr.io/YOUR_USERNAME/yarr:once-latest \
  --push .
```

Automated builds are also triggered by GitHub Actions when a version tag (`v*`) is pushed or via manual workflow dispatch.

See more:

* [Building from source code](doc/build.md)
* [Fever API support](doc/fever.md)

## credits

[Feather](http://feathericons.com/) for icons.
