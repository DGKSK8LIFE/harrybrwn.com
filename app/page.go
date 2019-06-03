package app

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// Page is a struct that represents a webpage
type Page struct {
	Title        string
	BodyTemplate string
	templates    []string
}

// Write will write the webpage to an io.Writer
func (p *Page) Write(w io.Writer) error {
	blob, err := newTemplate(p.allTemplates()...)
	if err != nil {
		return err
	}
	return execTemplate(w, blob, p)
}

// AddTemplate will add a template file to the page struct
func (p *Page) AddTemplate(files ...string) {
	p.templates = append(p.templates, files...)
}

func (p *Page) allTemplates() []string {
	return append(p.templates, p.BodyTemplate)
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
	// parseUserAgent(r.Header.Get("user-agent"))
	p.wrap.ServeHTTP(w, r)
}
