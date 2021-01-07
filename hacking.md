# hacking

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

## code of conduct

Be excellent to each other. Party on, dudes!
