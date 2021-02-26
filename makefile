VERSION=1.3
GITHASH=$(shell git rev-parse --short=8 HEAD)

CGO_ENABLED=1

GO_LDFLAGS  = -s -w
GO_LDFLAGS := $(GO_LDFLAGS) -X 'main.Version=$(VERSION)' -X 'main.GitHash=$(GITHASH)'

build_default:
	mkdir -p _output
	go build -tags "sqlite_foreign_keys release" -ldflags="$(GO_LDFLAGS)" -o _output/yarr src/main.go

build_macos:
	set GOOS=darwin
	set GOARCH=amd64
	mkdir -p _output/macos
	go build -tags "sqlite_foreign_keys release macos" -ldflags="$(GO_LDFLAGS)" -o _output/macos/yarr src/main.go
	cp platform/icon.png _output/macos/icon.png
	go run scripts/package_macos.go -outdir _output/macos -version "$(VERSION)"

build_linux:
	set GOOS=linux
	set GOARCH=386
	mkdir -p _output/linux
	go build -tags "sqlite_foreign_keys release linux" -ldflags="$(GO_LDFLAGS)" -o _output/linux/yarr src/main.go

build_windows:
	set GOOS=windows
	set GOARCH=386
	mkdir -p _output/windows
	go run scripts/generate_versioninfo.go -version "$(VERSION)" -outfile platform/versioninfo.rc
	windres -i platform/versioninfo.rc -O coff -o platform/versioninfo.syso
	go build -tags "sqlite_foreign_keys release windows" -ldflags="$(GO_LDFLAGS) -H windowsgui" -o _output/windows/yarr.exe main.go
