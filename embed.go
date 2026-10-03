// Package easytalk embeds the compiled frontend assets so EasyTalk can be
// distributed as a single binary with no external static files.
package easytalk

import (
	"embed"
	"io/fs"
)

//go:embed all:web/dist
var dist embed.FS

// FrontendFS returns the frontend assets rooted at the dist output directory.
func FrontendFS() (fs.FS, error) {
	return fs.Sub(dist, "web/dist")
}