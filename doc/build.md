Install `Go >= 1.14` and `gcc`. Get the source code:

    git clone https://github.com/nkanaev/yarr.git
    git clone https://github.com/nkanaev/gofeed.git yarr/gofeed

Then run one of the corresponding commands:

    # create a binary for the host os
    make build_macos    # -> _output/macos/yarr.app
    make build_linux    # -> _output/linux/yarr
    make build_windows  # -> _output/windows/yarr.exe

    # ... or start a dev server locally
    go run main.go      # starts a server at http://localhost:7070

    # ... or build a docker image
    docker build -t yarr .
