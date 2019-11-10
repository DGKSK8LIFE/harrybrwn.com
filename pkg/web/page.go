package web

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"path/filepath"

	"harrybrown.com/pkg/log"
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

	// Serve is a function used to serve http requests with the Page. If the
	// Serve field is not given the Page will execute it's internal template
	// blob and serve that.
	Serve func(w http.ResponseWriter, r *http.Request)

	// Data is an interface used as a vessel for getting data into the web
	// page template.
	Data interface{}

	// HotReload is a boolean that, if true, will cause the Page to call 'Init'
	// everytime it is written to an io.Writer. If there is no Serve function
	// set, then the page will be doing file IO upon every request.
	HotReload bool

	templates    []string
	blob         *template.Template
	baseTmplName string
}

// WriteTo will write the webpage to an io.Writer
func (p *Page) WriteTo(w io.Writer) (int64, error) {
	if p.HotReload {
		if err := p.Init(); err != nil {
			return 0, err
		}
	}

	b := &bytes.Buffer{}
	err := p.blob.ExecuteTemplate(b, p.baseTmplName, p)
	if err != nil {
		return 0, err
	}

	n, err := w.Write(b.Bytes())
	return int64(n), err
}

// ServerHTTP lets the Page struct impliment the http.Handler interface.
func (p *Page) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if p.Serve == nil {
		if _, err := p.WriteTo(w); err != nil {
			log.Println(err)
			NotFound(w, r)
		}
		return
	}
	p.Serve(w, r)
}

// AddTemplateFile will add a template file to the page struct
func (p *Page) AddTemplateFile(files ...string) {
	for _, file := range files {
		p.templates = append(p.templates, getfile(file))
	}
}

// Path returns the route path.
func (p *Page) Path() string {
	return p.RoutePath
}

// Handler returns the page that the method was called from.
func (p *Page) Handler() http.Handler {
	return p
}

// Init will initialize the Page by collecting all it's templates and
// creating the template blob.
func (p *Page) Init() (err error) {
	if len(p.templates) == 0 {
		p.templates = p.tmpls()
	}
	p.baseTmplName = BaseTemplateName
	p.blob, err = template.New(p.baseTmplName).ParseFiles(p.templates...)
	return err
}

func (p *Page) tmpls() []string {
	var (
		i      int
		t      string
		length = len(BaseTemplates) + 1
	)

	tmpls := make([]string, length)
	for i, t = range BaseTemplates {
		tmpls[i] = getfile(t)
	}
	tmpls[i+1] = getfile(p.Template)
	return tmpls
}

func (p *Page) templateCount() int {
	n := len(BaseTemplates) + len(p.templates)
	if len(p.Template) > 0 {
		n++
	}
	return n
}

var _ http.Handler = (*Page)(nil)
var _ Route = (*Page)(nil)
var _ io.WriterTo = (*Page)(nil)

func getfile(name string) string {
	if len(TemplateDir) < 1 {
		return name
	}
	return filepath.Join(TemplateDir, name)
}
