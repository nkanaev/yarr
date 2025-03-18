VERSION=2.4
GITHASH=$(shell git rev-parse --short=8 HEAD)

GO_TAGS    = sqlite_foreign_keys sqlite_json
GO_LDFLAGS = -s -w -X 'main.Version=$(VERSION)' -X 'main.GitHash=$(GITHASH)'

GO_FLAGS     = -tags "$(GO_TAGS)" -ldflags="$(GO_LDFLAGS)"
GO_FLAGS_GUI = -tags "$(GO_TAGS) gui" -ldflags="$(GO_LDFLAGS)"

export CGO_ENABLED=1

build_default:
	mkdir -p _output
	go build -tags "$(GO_TAGS)" -ldflags="$(GO_LDFLAGS)" -o _output/yarr ./cmd/yarr

build_macos:
	mkdir -p _output/macos
	GOOS=darwin go build -tags "$(GO_TAGS) macos" -ldflags="$(GO_LDFLAGS)" -o _output/macos/yarr ./cmd/yarr
	cp src/platform/icon.png _output/macos/icon.png
	go run ./cmd/package_macos -outdir _output/macos -version "$(VERSION)"

build_linux:
	mkdir -p _output/linux
	GOOS=linux go build -tags "$(GO_TAGS) linux" -ldflags="$(GO_LDFLAGS)" -o _output/linux/yarr ./cmd/yarr

build_windows:
	mkdir -p _output/windows
	go run ./cmd/generate_versioninfo -version "$(VERSION)" -outfile src/platform/versioninfo.rc
	windres -i src/platform/versioninfo.rc -O coff -o src/platform/versioninfo.syso
	GOOS=windows go build -tags "$(GO_TAGS) windows" -ldflags="$(GO_LDFLAGS) -H windowsgui" -o _output/windows/yarr.exe ./cmd/yarr

etc/icon.icns: etc/icon_macos.png
	mkdir -p etc/icon.iconset
	sips -s format png --resampleWidth 1024 etc/icon_macos.png --out etc/icon.iconset/icon_512x512@2x.png
	sips -s format png --resampleWidth  512 etc/icon_macos.png --out etc/icon.iconset/icon_512x512.png
	sips -s format png --resampleWidth  256 etc/icon_macos.png --out etc/icon.iconset/icon_256x256.png
	sips -s format png --resampleWidth  128 etc/icon_macos.png --out etc/icon.iconset/icon_128x128.png
	sips -s format png --resampleWidth   64 etc/icon_macos.png --out etc/icon.iconset/icon_32x32@2x.png
	sips -s format png --resampleWidth   32 etc/icon_macos.png --out etc/icon.iconset/icon_32x32.png
	sips -s format png --resampleWidth   16 etc/icon_macos.png --out etc/icon.iconset/icon_16x16.png
	iconutil -c icns etc/icon.iconset -o etc/icon.icns

darwin_arm64:
	GOOS=darwin GOARCH=arm64 go build $(GO_FLAGS) -o out/$@/yarr ./cmd/yarr

darwin_amd64:
	GOOS=darwin GOARCH=arm64 go build $(GO_FLAGS) -o out/$@/yarr ./cmd/yarr

linux_amd64:
	GOOS=linux GOARCH=amd64 go build $(GO_FLAGS) -o out/$@/yarr ./cmd/yarr

linux_arm64:
	GOOS=linux GOARCH=arm64 go build $(GO_FLAGS) -o out/$@/yarr ./cmd/yarr

darwin_arm64_gui: etc/icon.icns
	mkdir -p out/$@
	GOOS=darwin GOARCH=arm64 go build $(GO_FLAGS_GUI) -o out/$@/yarr ./cmd/yarr
	./etc/macos_package.sh $(VERSION) etc/icon.icns out/$@/yarr out/$@

darwin_amd64_gui: etc/icon.icns
	mkdir -p out/$@
	GOOS=darwin GOARCH=amd64 go build $(GO_FLAGS_GUI) -o out/$@/yarr ./cmd/yarr
	./etc/macos_package.sh $(VERSION) etc/icon.icns out/$@/yarr out/$@

serve:
	go run $(GO_FLAGS) ./cmd/yarr -db local.db

test:
	go test $(GO_FLAGS) ./...

.PHONY: serve test \
	darwin_amd64 darwin_amd64_gui \
	darwin_arm64 darwin_arm64_gui
