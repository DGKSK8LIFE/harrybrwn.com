package app

import (
	"fmt"
	"net/http"
	"os"

	"harrybrown.com/pkg/cmd"
	"harrybrown.com/pkg/log"
)

// RoutesCmd is the command that prints out the roues.
var RoutesCmd = cmd.Command{
	Syntax:      "routes",
	Description: "print out all the routes that the server is handling",
	Run: func() {
		for i, r := range Routes {
			fmt.Printf("%d: '%s'\n", i, r.Path())
		}
	},
}

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
