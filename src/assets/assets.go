package assets

import (
	"embed"
	"html/template"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
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

func Template(path string) *template.Template {
	var tmpl *template.Template
	tmpl, found := FS.templates[path]
	if !found {
		tmpl = template.Must(template.New(path).Delims("{%", "%}").Funcs(template.FuncMap{
			"inline": func(svg string) template.HTML {
				svgfile, err := FS.Open("graphicarts/" + svg)
				// should never happen
				if err != nil {
					log.Fatal(err)
				}
				defer svgfile.Close()

				content, err := ioutil.ReadAll(svgfile)
				// should never happen
				if err != nil {
					log.Fatal(err)
				}
				return template.HTML(content)
			},
		}).ParseFS(FS, path))
		if FS.embedded != nil {
			FS.templates[path] = tmpl
		}
	}
	return tmpl
}

func Render(path string, writer io.Writer, data interface{}) {
	tmpl := Template(path)
	tmpl.Execute(writer, data)
}

func init() {
	FS.templates = make(map[string]*template.Template)
}
