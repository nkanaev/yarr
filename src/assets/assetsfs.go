//go:build release
// +build release

package assets

import "embed"

//go:embed *.html
//go:embed graphicarts
//go:embed javascripts
//go:embed stylesheets
var embedded embed.FS

func init() {
	FS.embedded = &embedded
}
