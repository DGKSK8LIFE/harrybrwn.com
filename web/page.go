package web

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
)

var (
	// TemplateDir is a variable that can be set by the importer of this
	// library to as a prefix to any template files given to functions or
	// structs in this package.
	TemplateDir string

	// BaseTemplates is a variable that holds template file names that all
	// pages will use.
	BaseTemplates []string

	// BaseTemplateName is the name of the template that will be executed first
	// when a batch of templates are executed. (default "base")
	BaseTemplateName = "base"
)

// Page is a struct that represents a webpage
type Page struct {
	// the title of the web page
	Title string

	// Template is the main template used for the web page.
	Template string

	// RoutePath is the route used when serving the page.
	RoutePath string

	// Serve is a function used to serve http requests with the Page.
	Serve func(w http.ResponseWriter, r *http.Request)

	// Data is an interface used as a vessel for getting data into the web
	// page template.
	Data interface{}

	templates []string
	blob      *template.Template
	name      string
}

// Write will write the webpage to an io.Writer
func (p *Page) Write(w io.Writer) error {
	if err := p.init(); err != nil {
		return err
	}
	return p.blob.ExecuteTemplate(w, p.name, p)
}

// ServerHTTP lets the Page struct impliment the http.Handler interface.
func (p *Page) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if p.Serve == nil {
		if err := p.Write(w); err != nil {
			log.Println(err)
			NotFound(w, r)
		}
		return
	}
	p.Serve(w, r)
}

// AddTemplate will add a template file to the page struct
func (p *Page) AddTemplate(files ...string) {
	p.templates = append(p.templates, files...)
}

// Path returns the route path.
func (p *Page) Path() string {
	return p.RoutePath
}

// Handler returns the page that the method was called from.
func (p *Page) Handler() http.Handler {
	return p
}

func (p *Page) init() (err error) {
	// p.name = filepath.Base(p.Template)
	p.name = BaseTemplateName

	files := append(p.tmpls(), getfile(p.Template))
	p.blob, err = template.New(p.name).ParseFiles(files...)
	return err
}

func (p *Page) tmpls() (tmpls []string) {
	for _, t := range p.files() {
		tmpls = append(tmpls, getfile(t))
	}
	return tmpls
}

func (p *Page) files() (files []string) {
	files = BaseTemplates
	for _, f := range p.templates {
		files = append(files, f)
	}
	return files
}

var _ http.Handler = (*Page)(nil)
var _ Route = (*Page)(nil)

func getfile(name string) string {
	if len(TemplateDir) < 1 {
		return name
	}
	return filepath.Join(TemplateDir, name)
}
