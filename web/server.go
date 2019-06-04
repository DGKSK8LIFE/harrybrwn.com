package web

import (
	"fmt"
	"net/http"

	"harrybrown.com/email"
)

// Server is a server
type Server struct {
	mux    *http.ServeMux
	server *http.Server
	mailer *email.Sender
}

// NewServer creates a new server.
func NewServer() *Server {
	return &Server{mux: http.NewServeMux()}
}

// ListenAndServe will run the server.
func (s *Server) ListenAndServe(addr string) error {
	s.server = &http.Server{Addr: addr, Handler: s.mux}

	return s.server.ListenAndServe()
}

// HandlerHook is a hook that alows for the modification of handlers at
// runtime
var HandlerHook func(h http.Handler) http.Handler

// Handle registers the a path and a handler.
func (s *Server) Handle(path string, h http.Handler) {
	if HandlerHook != nil {
		h = HandlerHook(h)
	}
	s.mux.Handle(path, h)
}

// HandleFunc will register a new route with a HandlerFunc
func (s *Server) HandleFunc(path string, fn http.HandlerFunc) {
	s.mux.HandleFunc(path, fn)
}

// HandleRoutes will handle a list of routes.
func (s *Server) HandleRoutes(routes []Route) {
	for _, r := range routes {
		s.Handle(r.Path(), r.Handler())
	}
}

// HandleThing will handle a thing
func (s *Server) HandleThing(thing interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "thing\n")
	}
}

// Route defines an interface for various routes to be used in the Server struct.
type Route interface {
	Path() string
	Handler() http.Handler
}

// HTTPRoute is an http route
type HTTPRoute struct {
	RoutePath   string
	HTTPHandler http.Handler
}

// Path gets the path
func (r *HTTPRoute) Path() string {
	return r.RoutePath
}

// Handler gets the handler
func (r *HTTPRoute) Handler() http.Handler {
	return r.HTTPHandler
}

// NewRoute returns a basic Route interface
func NewRoute(path string, handler http.Handler) Route {
	return &HTTPRoute{RoutePath: path, HTTPHandler: handler}
}
