package frontend

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed all:build
var FsDir embed.FS

func BuildHTTPFS() http.FileSystem {
	build, err := fs.Sub(FsDir, "build")
	if err != nil {
		log.Fatal(err)
	}
	return http.FS(build)
}
