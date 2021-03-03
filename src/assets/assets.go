package assets

import (
	"embed"
	"html/template"
	"io"
	"io/ioutil"
	"io/fs"
	"os"
)

type assetsfs struct {
	embedded  *embed.FS
	templates map[string]*template.Template
}

var FS assetsfs

func (afs assetsfs) Open(name string) (fs.File, error) {
	if afs.embedded != nil {
		return afs.embedded.Open(name)
	}
	return os.DirFS("src/assets").Open(name)
}

func Render(path string, writer io.Writer, data interface{}) {
	var tmpl *template.Template
	tmpl, found := FS.templates[path]
	if !found {
		tmpl = template.Must(template.New(path).Delims("{%", "%}").Funcs(template.FuncMap{
			"inline": func(svg string) template.HTML {
				svgfile, _ := FS.Open("graphicarts/" + svg)
				content, _ := ioutil.ReadAll(svgfile)
				svgfile.Close()
				return template.HTML(content)
			},
		}).ParseFS(FS, path))
		if FS.embedded != nil {
			FS.templates[path] = tmpl
		}
	}
	tmpl.Execute(writer, data)
}

func init() {
	FS.templates = make(map[string]*template.Template)
}
