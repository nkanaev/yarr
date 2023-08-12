## Compilation

Install `Go >= 1.17` and `GCC`. Get the source code:

    git clone https://github.com/nkanaev/yarr.git

Then run one of the corresponding commands:

    # create an executable for the host os
    make build_macos    # -> _output/macos/yarr.app
    make build_linux    # -> _output/linux/yarr
    make build_windows  # -> _output/windows/yarr.exe

    # host-specific cli version (no gui)
    make build_default  # -> _output/yarr

    # ... or start a dev server locally
    make serve          # starts a server at http://localhost:7070

    # ... or build a docker image
    docker build -t yarr .

## ARM compilation

The instructions below are to cross-compile *yarr* to `Linux/ARM*`.

Build:

    docker build -t yarr.arm -f dockerfile.arm .

Test:

    # inside host
    docker run -it --rm yarr.arm

    # then, inside container
    cd /root/out
    qemu-aarch64 -L /usr/aarch64-linux-gnu/ yarr.arm64

Extract files from images:

    CID=$(docker create yarr.arm)
    docker cp -a "$CID:/root/out" .
    docker rm "$CID"
