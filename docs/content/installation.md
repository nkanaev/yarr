---
title: "Installation"
---

Prebuilt binaries for Linux, MacOS and Windows are available on the
[releases page](https://github.com/nkanaev/yarr/releases/latest).
Archives follow the naming convention `yarr_{OS}_{ARCH}[_gui].zip`:

* `OS` is the target operating system
* `ARCH` is the CPU architecture (`arm64` for AArch64, `amd64` for X86-64)
* `-gui` indicates that the binary ships with the GUI (tray icon); omit it for a command line application

## MacOS

Place `yarr.app` in the `/Applications` folder,
[open the app](https://support.apple.com/en-gb/guide/mac-help/mh40616/mac),
click the anchor menu bar icon and select "Open".

## Windows

Open `yarr.exe`, click the anchor system tray icon and select "Open".

## Linux

Place `yarr` in `$HOME/.local/bin` and run the
[install script](https://raw.githubusercontent.com/nkanaev/yarr/master/etc/install-linux.sh).

## Self-hosting

For self-hosting, see `yarr -h` for auth, tls & server configuration flags.

## Building from source

Prerequisites:

* Go >= 1.23
* C Compiler (GCC / Clang / ...)
* Zig >= 0.14.0 (optional, for cross-compiling CLI versions)
* binutils (optional, for building Windows GUI version)

Get the source code:

    git clone https://github.com/nkanaev/yarr.git

Compile:

    # create cli for the host OS/architecture
    make host               # out/yarr

    # create GUI, works only in the target OS
    make windows_amd64_gui  # out/windows_amd64_gui/yarr.exe
    make windows_arm64_gui  # out/windows_arm64_gui/yarr.exe
    make darwin_arm64_gui   # out/darwin_arm64_gui/yarr.app
    make darwin_amd64_gui   # out/darwin_amd64_gui/yarr.app

    # create cli, cross-compiles within any OS/architecture
    make linux_amd64
    make linux_arm64
    make linux_armv7
    make windows_amd64
    make windows_arm64

    # ... or build a docker image
    docker build -t yarr -f etc/dockerfile .

### ARM compilation

To cross-compile *yarr* to `Linux/ARM*`:

    docker build -t yarr.arm -f etc/dockerfile.arm .
