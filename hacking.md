# hacking

If you have any questions/suggestions/proposals,
you can reach out the author via e-mail (nkanaev@live.com)
or mastodon (https://fosstodon.org/@nkanaev).

## build

Install `Go >= 1.14` and `gcc`. Get the source code:

```sh
git clone https://github.com/nkanaev/yarr.git
git clone https://github.com/nkanaev/gofeed.git
mv gofeed yarr
cd yarr
```

Then:

```sh
# create a binary for the host os
make build_macos    # -> _output/macos/yarr.app
make build_linux    # -> _output/linux/yarr
make build_windows  # -> _output/windows/yarr.exe

# ... or run locally (for testing & hacking)
go run main.go      # starts a server at http://localhost:7070
```

## plans

- test across 3 platforms (macos, linux, windows)
- prebuilt binaries
- GUI-less mode (no tray icon)
- feeds health checker
- mobile & tablet layout
- parameters (`--[no]-gui`, `--addr`, ...)
- Fever API support
- keyboard navigation

## code of conduct

Be excellent to each other. Party on, dudes!
