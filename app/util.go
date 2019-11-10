package app

import (
	"net/http"
	"os"

	"harrybrown.com/pkg/log"
	"harrybrown.com/pkg/web"
)

func init() {
	// log.DefaultLogger = log.NewPlainLogger(os.Stdout)
	web.HandlerHook = NewLogger
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
