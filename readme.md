# yarr

**yarr** (yet another rss reader) is a web-based feed aggregator which can be used both
as a desktop application and a personal self-hosted server.

It is written in Go with the frontend in Vue.js. The storage is backed by SQLite.

![screenshot](etc/promo.png)

## Usage

The latest prebuilt binaries for Linux/MacOS/Windows are available
[here](https://github.com/nkanaev/yarr/releases/latest).

### MacOS

Download `yarr-*-macos-amd64.zip`, unzip it, place `yarr.app` in `/Applications` folder, [open the app][macos-open], click the anchor menu bar icon, select "Open".

[macos-open]: https://support.apple.com/en-gb/guide/mac-help/mh40616/mac

### Windows

Download `yarr-*-windows-amd64.zip`, unzip it, open `yarr.exe`, click the anchor system tray icon, select "Open".

### Linux

Download `yarr-*-linux-amd64.zip`, unzip it, place `yarr` in `$HOME/.local/bin`
and run [the provided script](etc/install-linux.sh).

For self-hosting, see `yarr -h` for auth, tls & server configuration flags.
For building from source code, see [BUILD.md](BUILD.md)

## Credits

[Feather](http://feathericons.com/) for icons.
