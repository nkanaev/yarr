VERSION=2.4
GITHASH=$(shell git rev-parse --short=8 HEAD)

GO_TAGS    = sqlite_foreign_keys sqlite_json
GO_LDFLAGS = -s -w -X 'main.Version=$(VERSION)' -X 'main.GitHash=$(GITHASH)'

GO_FLAGS         = -tags "$(GO_TAGS)"     -ldflags="$(GO_LDFLAGS)"
GO_FLAGS_GUI     = -tags "$(GO_TAGS) gui" -ldflags="$(GO_LDFLAGS)"
GO_FLAGS_GUI_WIN = -tags "$(GO_TAGS) gui" -ldflags="$(GO_LDFLAGS) -H windowsgui"

export CGO_ENABLED=1

default: test host

# platform-specific files

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

src/platform/versioninfo.rc:
	./etc/windows_versioninfo.sh -version "$(VERSION)" -outfile src/platform/versioninfo.rc
	windres -i src/platform/versioninfo.rc -O coff -o src/platform/versioninfo.syso

# build targets

host:
	go build $(GO_FLAGS) -o out/yarr ./cmd/yarr

darwin_amd64:
	# not supported yet
	# CC="zig cc -target x86_64-macos-none" GOOS=darwin GOARCH=arm64 go build $(subst -s ,,$(GO_FLAGS)) -o out/$@/yarr ./cmd/yarr

darwin_arm64:
	# not supported yet
	# CC="zig cc -target aarch64-macos-none" GOOS=darwin GOARCH=arm64 go build $(subst -s ,,$(GO_FLAGS)) -o out/$@/yarr ./cmd/yarr

linux_amd64:
	CC="zig cc -target x86_64-linux-musl -O2 -g0" CGO_CFLAGS="-D_LARGEFILE64_SOURCE" GOOS=linux GOARCH=amd64 \
	go build $(GO_FLAGS) -o out/$@/yarr ./cmd/yarr

linux_arm64:
	CC="zig cc -target aarch64-linux-musl -O2 -g0" CGO_CFLAGS="-D_LARGEFILE64_SOURCE" GOOS=linux GOARCH=arm64 \
	go build $(GO_FLAGS) -o out/$@/yarr ./cmd/yarr

linux_armv7:
	CC="zig cc -target arm-linux-musleabihf -O2 -g0" CGO_CFLAGS="-D_LARGEFILE64_SOURCE" GOOS=linux GOARCH=arm GOARM=7 \
	go build $(GO_FLAGS) -o out/$@/yarr ./cmd/yarr

windows_amd64:
	CC="zig cc -target x86_64-windows-gnu" GOOS=windows GOARCH=amd64 go build $(GO_FLAGS) -o out/$@/yarr ./cmd/yarr

windows_arm64:
	CC="zig cc -target aarch64-windows-gnu" GOOS=windows GOARCH=arm64 go build $(GO_FLAGS) -o out/$@/yarr ./cmd/yarr

darwin_arm64_gui: etc/icon.icns
	GOOS=darwin GOARCH=arm64 go build $(GO_FLAGS_GUI) -o out/$@/yarr ./cmd/yarr
	./etc/macos_package.sh $(VERSION) etc/icon.icns out/$@/yarr out/$@

darwin_amd64_gui: etc/icon.icns
	GOOS=darwin GOARCH=amd64 go build $(GO_FLAGS_GUI) -o out/$@/yarr ./cmd/yarr
	./etc/macos_package.sh $(VERSION) etc/icon.icns out/$@/yarr out/$@

windows_amd64_gui: windows_versioninfo
	GOOS=windows GOARCH=amd64 go build $(GO_FLAGS_GUI_WIN) -o out/$@/yarr ./cmd/yarr

windows_arm64_gui: src/platform/versioninfo.rc
	GOOS=windows GOARCH=arm64 go build $(GO_FLAGS_GUI_WIN) -o out/$@/yarr ./cmd/yarr

serve:
	go run $(GO_FLAGS) ./cmd/yarr -db local.db

test:
	go test $(GO_FLAGS) ./...

.PHONY: \
	host \
	darwin_amd64 darwin_amd64_gui \
	darwin_arm64 darwin_arm64_gui \
	windows_amd64 windows_amd64_gui \
	windows_arm64 windows_arm64_gui \
	serve test
