//go:build release
// +build release

package assets

import "embed"

//go:embed *.html
//go:embed graphicarts
//go:embed javascripts
//go:embed stylesheets
//go:embed manifest.json
var embedded embed.FS

func init() {
	FS.embedded = &embedded
}
