# yarr

**yarr** (yet another rss reader) is a web-based feed aggregator which can be used both
as a desktop application and a self-hosted server.

It is written in Go with the frontend in Vue.js. The storage is backed by SQLite.

![screenshot](etc/promo.png)

[download](https://github.com/nkanaev/yarr/releases/latest)

## build

Install `Go >= 1.16` and `gcc`. Get the source code:

    git clone --recurse-submodules https://github.com/nkanaev/yarr.git

Then run one of the corresponding commands:

    # create an executable for the host os
    make build_macos    # -> _output/macos/yarr.app
    make build_linux    # -> _output/linux/yarr
    make build_windows  # -> _output/windows/yarr.exe

    # ... or start a dev server locally
    go run main.go      # starts a server at http://localhost:7070

    # ... or build a docker image
    docker build -t yarr .

## credits

[Feather](http://feathericons.com/) for icons.
