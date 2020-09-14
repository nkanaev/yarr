# yarr (beta)

yet another rss reader.

![screenshot](https://github.com/nkanaev/yarr/blob/master/artwork/promo.png?raw=true)

*yarr* is a server written in Go with the frontend in Vue.js. The storage is backed by SQLite.

The goal of the project is to provide a desktop application accessible via web browser.
Longer-term plans include a self-hosted solution for individuals.

There are plans to add support for mobile & table resolutions.
Support for 3rd-party applications (via Fever API) is being considered.

## build

Install `Go >= 1.14` and `gcc`, then run:

```sh
$ git clone https://github.com/nkanaev/yarr.git
$ git clone https://github.com/nkanaev/gofeed.git
$ mv gofeed yarr
$ cd yarr && make build_macos
```

## plans

- test across 3 platforms (macos, linux, windows)
- binaries
- gui-less mode (no tray icon)
- feeds health checker
- mobile & tablet layout
- parameters (`--[no]-gui`, `--addr`, ...)
- Fever API support
- keyboard navigation
