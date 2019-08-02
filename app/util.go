package app

import (
	"net/http"
	"os"

	"harrybrown.com/pkg/log"
)

// NewLogger creates a new logger that will intercept a handler and replace it
// with one that has logging functionality.
func NewLogger(h http.Handler) http.Handler {
	return &pageLogger{
		wrap: h,
		// l:    log.New(os.Stdout, "", log.LstdFlags),
		l: log.NewColorLogger(os.Stdout, log.Purple),
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
