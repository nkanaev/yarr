ASSETS = assets/javascripts/* assets/stylesheets/* assets/graphicarts/* assets/index.html
CGO_ENABLED=1

default: bundle

server/assets_bundle.go: $(ASSETS)
	go run scripts/bundle_assets.go >/dev/null

bundle: server/assets_bundle.go

build_macos: bundle
	set GOOS=darwin
	set GOARCH=amd64
	mkdir -p _output/macos
	go build -tags "sqlite_foreign_keys release macos" -ldflags="-s -w" -o _output/macos/yarr main.go
	cp artwork/icon.png _output/macos/icon.png
	go run scripts/package_macos.go _output/macos

build_linux: bundle
	set GOOS=linux
	set GOARCH=386
	mkdir -p _output/linux
	go build -tags "sqlite_foreign_keys release linux" -ldflags="-s -w" -o _output/linux/yarr main.go

build_windows: bundle
	set GOOS=windows
	set GOARCH=386
	mkdir -p _output/windows
	go build -tags "sqlite_foreign_keys release windows" -ldflags="-s -w -H windowsgui" -o _output/windows/yarr.exe main.go
