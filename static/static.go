package static

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:generate npm run build
//go:embed all:dist
var dist embed.FS

func Assets() (http.Handler, error) {
	a, err := fs.Sub(dist, "dist")
	if err != nil {
		return nil, err
	}
	fs := http.FileServer(http.FS(a))
	return http.StripPrefix("/", fs), nil
}
