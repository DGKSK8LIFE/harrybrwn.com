package app

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// FileServer is a custum file server.
type FileServer struct {
	fserver http.Handler
}

// NewFileServer will create a new FileServer
func NewFileServer(dir string) http.Handler {
	return &FileServer{http.FileServer(http.Dir(dir))}
}

func (fs *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	size := r.URL.Query().Get("size")
	if len(size) > 0 {
		r.URL.Path = pointToImage(r.URL, size)
	}

	fs.fserver.ServeHTTP(w, r)
}

func pointToImage(u *url.URL, size string) string {
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

// UserAgent is a struct that represents a user-agent
type UserAgent struct {
	OS      *OSInfo
	Browser string
}

// OSInfo is a grouping of information about a computers operating system.
type OSInfo struct {
	Name    string
	Version float32
	Arch    string
}

func parseUserAgent(agent string) {
	parts := strings.Split(agent, " ")
	fmt.Println("User-Agent:", parts)
}

// NewLogger creates a new logger that will intercept a handler and replace it
// with one that has logging functionality.
func NewLogger(h http.Handler) http.Handler {
	return &pageLogger{wrap: h}
}

type pageLogger struct {
	wrap http.Handler
}

func (p *pageLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.Method, r.Proto, r.URL)
	p.wrap.ServeHTTP(w, r)
}
