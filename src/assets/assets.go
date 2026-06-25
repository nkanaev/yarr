//go:build debug

package assets

import (
	"html/template"
	"io/fs"
	"os"
)

const rootDir = "src/assets"

func Templates() *template.Template {
	return template.Must(template.ParseGlob(rootDir + "/templates/*.html"))
}

func StaticFS() fs.FS {
	return os.DirFS(rootDir)
}
