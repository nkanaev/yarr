VERSION=2.0
GITHASH=$(shell git rev-parse --short=8 HEAD)

CGO_ENABLED=1

GO_LDFLAGS  = -s -w
GO_LDFLAGS := $(GO_LDFLAGS) -X 'main.Version=$(VERSION)' -X 'main.GitHash=$(GITHASH)'

build_default:
	mkdir -p _output
	go build -tags "sqlite_foreign_keys release" -ldflags="$(GO_LDFLAGS)" -o _output/yarr src/main.go

build_macos:
	mkdir -p _output/macos
	GOOS=darwin GOARCH=amd64 go build -tags "sqlite_foreign_keys release macos" -ldflags="$(GO_LDFLAGS)" -o _output/macos/yarr src/main.go
	cp src/platform/icon.png _output/macos/icon.png
	go run bin/package_macos.go -outdir _output/macos -version "$(VERSION)"

build_linux:
	mkdir -p _output/linux
	GOOS=linux GOARCH=amd64 go build -tags "sqlite_foreign_keys release linux" -ldflags="$(GO_LDFLAGS)" -o _output/linux/yarr src/main.go

build_windows:
	mkdir -p _output/windows
	go run bin/generate_versioninfo.go -version "$(VERSION)" -outfile src/platform/versioninfo.rc
	windres -i src/platform/versioninfo.rc -O coff -o src/platform/versioninfo.syso
	GOOS=windows GOARCH=amd64 go build -tags "sqlite_foreign_keys release windows" -ldflags="$(GO_LDFLAGS) -H windowsgui" -o _output/windows/yarr.exe src/main.go

serve:
	go run -tags "sqlite_foreign_keys" src/main.go -db local.db

test:
	cd src && go test -tags "sqlite_foreign_keys release" ./...
