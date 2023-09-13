# yarr

**yarr** (yet another rss reader) is a web-based feed aggregator which can be used both
as a desktop application and a personal self-hosted server.

The app is a single binary with an embedded database (SQLite).

![screenshot](etc/promo.png)

## usage

The latest prebuilt binaries for Linux/MacOS/Windows AMD64 are available
[here](https://github.com/nkanaev/yarr/releases/latest).

### macos

Download `yarr-*-macos64.zip`, unzip it, place `yarr.app` in `/Applications` folder, [open the app][macos-open], click the anchor menu bar icon, select "Open".

[macos-open]: https://support.apple.com/en-gb/guide/mac-help/mh40616/mac

### windows

Download `yarr-*-windows64.zip`, unzip it, open `yarr.exe`, click the anchor system tray icon, select "Open".

### linux

Download `yarr-*-linux64.zip`, unzip it, place `yarr` in `$HOME/.local/bin`
and run [the script](etc/install-linux.sh).

For self-hosting, see `yarr -h` for auth, tls & server configuration flags.
For building from source code, see [build.md](build.md)

## Fever API support

Fever API is a kind of RSS HTTP API interface, because the Fever API definition is not very clear, so the implementation of Fever server and Client may have some compatibility problems.

The Fever API implemented by Yarr is based on the Fever API spec: https://github.com/DigitalDJ/tinytinyrss-fever-plugin/blob/master/fever-api.md.

Here are some Apps that have been tested to work with yarr.  Feel free to test other Clients/Apps and update the list here.

>  Different apps support different URL/Address formats.  Please note whether the URL entered has `http://` scheme and `/` suffix.

| App                                                                       | Platforms        | Config Server URL                                   |
|:------------------------------------------------------------------------- | ---------------- |:--------------------------------------------------- |
| [Reeder](https://reederapp.com/)                                          | MacOS<br>iOS     | 127.0.0.1:7070/fever<br>http://127.0.0.1:7070/fever |
| [ReadKit](https://readkit.app/)                                           | MacOS<br>iOS     | http://127.0.0.1:7070/fever                         |
| [Fluent Reader](https://github.com/yang991178/fluent-reader)              | MacOS<br>Windows | http://127.0.0.1:7070/fever/                        |
| [Unread](https://apps.apple.com/us/app/unread-an-rss-reader/id1363637349) | iOS              | http://127.0.0.1:7070/fever                         |
| [Fiery Feeds](https://voidstern.net/fiery-feeds)                          | MacOS<br>iOS     | http://127.0.0.1:7070/fever                         |


If you are having trouble using Fever, please open an issue and @icefed, thanks.

## credits

[Feather](http://feathericons.com/) for icons.
