ASSETS = assets/javascripts/* assets/stylesheets/* assets/graphicarts/* assets/index.html

CGO_ENABLED=1

default: build

server/assets_bundle.go: $(ASSETS)
	go run bundle.go >/dev/null

bundle: server/assets_bundle.go

build: build_mac build_nix build_win

build_mac: bundle
	set GOOS=darwin
	set GOARCH=amd64
	mkdir -p build/mac
	go build -tags "sqlite_foreign_keys release mac" -ldflags="-s -w" -o build/mac/yarr main.go

build_nix: bundle
	set GOOS=linux
	set GOARCH=386
	mkdir -p build/nix
	go build -tags "sqlite_foreign_keys release nix" -ldflags="-s -w" -o build/nix/yarr main.go

build_win: bundle
	set GOOS=windows
	set GOARCH=386
	mkdir -p build/win
	go build -tags "sqlite_foreign_keys release win" -ldflags="-s -w -H windowsgui" -o build/win/yarr.exe main.go

.PHONY: default bundle build build_mac build_nix build_win
