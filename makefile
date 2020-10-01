VERSION=v1.0
GITHASH=$(shell git rev-parse --short=8 HEAD)

ASSETS = assets/javascripts/* assets/stylesheets/* assets/graphicarts/* assets/index.html
CGO_ENABLED=1

GO_LDFLAGS  = -s -w
GO_LDFLAGS := $(GO_LDFLAGS) -X 'main.Version=$(VERSION)' -X 'main.GitHash=$(GITHASH)'

default: bundle

server/assets.go: $(ASSETS)
	go run scripts/bundle_assets.go >/dev/null

bundle: server/assets.go

build_macos: bundle
	set GOOS=darwin
	set GOARCH=amd64
	mkdir -p _output/macos
	go build -tags "sqlite_foreign_keys release macos" -ldflags="$(GO_LDFLAGS)" -o _output/macos/yarr main.go
	cp artwork/icon.png _output/macos/icon.png
	go run scripts/package_macos.go _output/macos

build_linux: bundle
	set GOOS=linux
	set GOARCH=386
	mkdir -p _output/linux
	go build -tags "sqlite_foreign_keys release linux" -ldflags="$(GO_LDFLAGS)" -o _output/linux/yarr main.go

build_windows: bundle
	set GOOS=windows
	set GOARCH=386
	mkdir -p _output/windows
	windres -i artwork/versioninfo.rc -O coff -o platform/versioninfo.syso
	go build -tags "sqlite_foreign_keys release windows" -ldflags="$(GO_LDFLAGS) -H windowsgui" -o _output/windows/yarr.exe main.go
