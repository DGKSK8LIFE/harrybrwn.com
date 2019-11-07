package app

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
)

// FileServer is a custum file server.
type FileServer struct {
	fserver http.Handler
}

// NewFileServer will create a new FileServer
func NewFileServer(dir string) http.Handler {
	return &FileServer{
		http.StripPrefix(fmt.Sprintf("/%s/", dir), http.FileServer(http.Dir(dir))),
	}
}

func (fs *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	size := r.URL.Query().Get("size")

	if len(size) > 0 {
		r.URL.Path = findImage(r.URL, size)
	}

	fs.fserver.ServeHTTP(w, r)
}

func findImage(u *url.URL, size string) string {
	var dir string
	switch size {
	case "xs":
		dir = "563x750"
	case "sm":
		dir = "1125x1500"
	case "md":
		dir = "1688x2251"
	case "lg":
		dir = "2250x3000"
	}
	folder, file := path.Split(u.Path)
	return path.Join(folder, dir, file)
}
