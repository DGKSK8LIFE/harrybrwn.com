package app

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"harrybrown.com/pkg/log"
	"harrybrown.com/pkg/web"
)

func init() {
	// log.DefaultLogger = log.NewPlainLogger(os.Stdout)
	web.HandlerHook = NewLogger
	web.DefaultErrorHandler = http.HandlerFunc(NotFound)
}

// NewLogger creates a new logger that will intercept a handler and replace it
// with one that has logging functionality.
func NewLogger(h http.Handler) http.Handler {
	return &pageLogger{
		wrap: h,
		l:    log.NewPlainLogger(os.Stdout),
	}
}

type pageLogger struct {
	wrap http.Handler
	l    log.PrintLogger
}

func (p *pageLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.l.Printf("%s %s %s\n", r.Method, r.Proto, r.URL)
	p.wrap.ServeHTTP(w, r)
}

// NotFound handles requests that generate a 404 error
func NotFound(w http.ResponseWriter, r *http.Request) {
	var tmplNotFound = template.Must(template.ParseFiles(
		"templates/pages/404.html",
		"templates/index.html",
	))
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if err := tmplNotFound.ExecuteTemplate(w, "base", nil); err != nil {
		fmt.Fprintf(w, "%s", err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}
