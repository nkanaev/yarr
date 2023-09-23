VERSION=2.4
GITHASH=$(shell git rev-parse --short=8 HEAD)

CGO_ENABLED=1

GO_LDFLAGS  = -s -w
GO_LDFLAGS := $(GO_LDFLAGS) -X 'main.Version=$(VERSION)' -X 'main.GitHash=$(GITHASH)'

build_default:
	mkdir -p _output
	go build -tags "sqlite_foreign_keys" -ldflags="$(GO_LDFLAGS)" -o _output/yarr ./cmd/yarr

build_macos:
	mkdir -p _output/macos
	GOOS=darwin GOARCH=amd64 go build -tags "sqlite_foreign_keys macos" -ldflags="$(GO_LDFLAGS)" -o _output/macos/yarr ./cmd/yarr
	cp src/platform/icon.png _output/macos/icon.png
	go run ./cmd/package_macos -outdir _output/macos -version "$(VERSION)"

build_linux:
	mkdir -p _output/linux
	GOOS=linux GOARCH=amd64 go build -tags "sqlite_foreign_keys linux" -ldflags="$(GO_LDFLAGS)" -o _output/linux/yarr ./cmd/yarr

build_windows:
	mkdir -p _output/windows
	go run ./cmd/generate_versioninfo -version "$(VERSION)" -outfile src/platform/versioninfo.rc
	windres -i src/platform/versioninfo.rc -O coff -o src/platform/versioninfo.syso
	GOOS=windows GOARCH=amd64 go build -tags "sqlite_foreign_keys windows" -ldflags="$(GO_LDFLAGS) -H windowsgui" -o _output/windows/yarr.exe ./cmd/yarr

serve:
	go run -tags "sqlite_foreign_keys" ./cmd/yarr -db local.db

test:
	cd src && go test -tags "sqlite_foreign_keys" ./...
