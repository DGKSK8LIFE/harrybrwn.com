package app

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

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
	// parseUserAgent(r.Header.Get("user-agent"))
	p.wrap.ServeHTTP(w, r)
}
