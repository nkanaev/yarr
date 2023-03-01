@echo off
set VERSION=2.3

set CGO_ENABLED=0

set GO_LDFLAGS=-s -w
set GO_LDFLAGS=%GO_LDFLAGS% -X "main.Version=%VERSION%" -X "main.GitHash=95ebbb9d"

rmdir /s /q _output\windows
mkdir _output\windows

go run bin\generate_versioninfo.go -version "%VERSION%" -outfile src\platform\versioninfo.rc
windres -i src\platform\versioninfo.rc -O coff -o src\platform\versioninfo.syso

set GOOS=windows
set GOARCH=amd64
go build -tags "sqlite_foreign_keys release windows" -ldflags="%GO_LDFLAGS% -H windowsgui" -o _output\windows\yarr.exe src\main.go

echo Build completed!
