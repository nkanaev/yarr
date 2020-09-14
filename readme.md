# yarr (beta)

yet another rss reader.

![screenshot](https://github.com/nkanaev/yarr/blob/master/artwork/promo.png?raw=true)

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
