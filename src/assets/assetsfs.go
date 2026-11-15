//go:build !debug

package assets

import (
	"embed"
	"html/template"
	"io/fs"
)

//go:embed bundle/*
var static embed.FS

//go:embed templates/*.html
var templates embed.FS

func Templates() *template.Template {
	return template.Must(template.ParseFS(templates, "templates/*.html"))
}

func StaticFS() fs.FS {
	fs, _ := fs.Sub(static, "bundle")
	return fs
}
