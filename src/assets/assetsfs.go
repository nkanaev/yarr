//go:build !debug

package assets

import (
	"embed"
	"html/template"
	"io/fs"
)

//go:embed static
var static embed.FS

//go:embed templates
var templates embed.FS

func Templates() *template.Template {
    return template.Must(template.ParseFS(templates, "*.html"))
}

func StaticFS() fs.FS {
    return static
}
