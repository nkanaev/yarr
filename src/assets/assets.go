package assets

import (
	"embed"
	"html/template"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type assetsfs struct {
	embedded  *embed.FS
	templates map[string]*template.Template
}

var FS assetsfs

var partialTmpl *template.Template

func (afs assetsfs) Open(name string) (fs.File, error) {
	if afs.embedded != nil {
		return afs.embedded.Open(name)
	}
	return os.DirFS("src/assets").Open(name)
}

func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"inline": func(svg string) template.HTML {
			svgfile, err := FS.Open("graphicarts/" + svg)
			if err != nil {
				log.Fatal(err)
			}
			defer svgfile.Close()

			content, err := ioutil.ReadAll(svgfile)
			if err != nil {
				log.Fatal(err)
			}
			return template.HTML(content)
		},
	}
}

func Template(path string) *template.Template {
	var tmpl *template.Template
	tmpl, found := FS.templates[path]
	if !found {
		// Parse the requested template along with all partial templates
		// so that {% template "partial_name" .data %} works
		patterns := []string{path}

		// Check if templates directory exists and add partial templates
		if dir, err := fs.ReadDir(FS, "templates"); err == nil {
			for _, entry := range dir {
				if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".html") {
					patterns = append(patterns, "templates/"+entry.Name())
				}
			}
		}

		tmpl = template.Must(template.New(path).Delims("{%", "%}").Funcs(templateFuncs()).ParseFS(FS, patterns...))
		if FS.embedded != nil {
			FS.templates[path] = tmpl
		}
	}
	return tmpl
}

// PartialTemplate returns a template set containing all partial templates.
// Templates are defined using {% define "name" %} blocks and can be executed
// individually via ExecuteTemplate(w, "name", data).
func PartialTemplate() *template.Template {
	if partialTmpl != nil && FS.embedded != nil {
		return partialTmpl
	}

	// Collect all template files from templates/ directory
	tmpl := template.New("partials").Delims("{%", "%}").Funcs(templateFuncs())

	// Read template files
	dir, err := fs.ReadDir(FS, "templates")
	if err != nil {
		log.Printf("Warning: could not read templates directory: %v", err)
		return tmpl
	}

	for _, entry := range dir {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".html") {
			continue
		}

		f, err := FS.Open("templates/" + entry.Name())
		if err != nil {
			log.Printf("Warning: could not open template %s: %v", entry.Name(), err)
			continue
		}

		content, err := ioutil.ReadAll(f)
		f.Close()
		if err != nil {
			log.Printf("Warning: could not read template %s: %v", entry.Name(), err)
			continue
		}

		_, err = tmpl.Parse(string(content))
		if err != nil {
			log.Printf("Warning: could not parse template %s: %v", entry.Name(), err)
			continue
		}
	}

	if FS.embedded != nil {
		partialTmpl = tmpl
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
