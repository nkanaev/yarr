VERSION := 2.3
GITHASH := $(shell git rev-parse --short=8 HEAD)
OUTPUT_DIR := _output

CGO_ENABLED := 1
GO_LDFLAGS  := -s -w
GO_LDFLAGS := $(GO_LDFLAGS) -X 'main.Version=$(VERSION)' -X 'main.GitHash=$(GITHASH)'

build_default: deps
	mkdir -p _output
	go build -tags "sqlite_foreign_keys release" -ldflags="$(GO_LDFLAGS)" -o $(OUTPUT_DIR)/yarr src/main.go

build_macos: deps
	mkdir -p $(OUTPUT_DIR)/macos
	GOOS=darwin GOARCH=amd64 go build -tags "sqlite_foreign_keys release macos" -ldflags="$(GO_LDFLAGS)" -o $(OUTPUT_DIR)/macos/yarr src/main.go
	cp src/platform/icon.png _output/macos/icon.png
	go run bin/package_macos.go -outdir _output/macos -version "$(VERSION)"

build_linux: deps
	mkdir -p $(OUTPUT_DIR)/linux
	GOOS=linux GOARCH=amd64 go build -tags "sqlite_foreign_keys release linux" -ldflags="$(GO_LDFLAGS)" -o $(OUTPUT_DIR)/linux/yarr src/main.go

build_windows: deps
	mkdir -p $(OUTPUT_DIR)/windows
	go run bin/generate_versioninfo.go -version "$(VERSION)" -outfile src/platform/versioninfo.rc
	windres -i src/platform/versioninfo.rc -O coff -o src/platform/versioninfo.syso
	GOOS=windows GOARCH=amd64 go build -tags "sqlite_foreign_keys release windows" -ldflags="$(GO_LDFLAGS) -H windowsgui" -o $(OUTPUT_DIR)/windows/yarr.exe src/main.go

clean:
	rm -rf _output

deps:
	go mod download

serve: deps
	go run -tags "sqlite_foreign_keys" src/main.go -db local.db

test: deps
	cd src && go test -tags "sqlite_foreign_keys release" ./...

fmt:
	cd src && go fmt ./...
