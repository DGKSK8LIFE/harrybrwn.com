package web

import (
	"fmt"
	"net/http"
)

// Router is an http router.
type Router struct {
	mux    *http.ServeMux
	server *http.Server
}

// NewRouter creates a new router.
func NewRouter() *Router {
	return &Router{mux: http.NewServeMux()}
}

// ListenAndServe will run the server.
func (s *Router) ListenAndServe(addr string) error {
	s.server = &http.Server{Addr: addr, Handler: s.mux}
	if HandlerHook != nil {
		s.server.Handler = HandlerHook(s.mux)
	}

	return s.server.ListenAndServe()
}

// HandlerHook is a hook that alows for the modification of handlers at
// runtime
var HandlerHook func(h http.Handler) http.Handler

// Handle registers the a path and a handler.
func (s *Router) Handle(path string, h http.Handler) {
	s.mux.Handle(path, h)
}

// HandleFunc will register a new route with a HandlerFunc
func (s *Router) HandleFunc(path string, fn http.HandlerFunc) {
	s.Handle(path, http.HandlerFunc(fn))
}

// HandleRoutes will handle a list of routes.
func (s *Router) HandleRoutes(routes []Route) {
	for _, r := range routes {
		s.Handle(r.Path(), r.Handler())
	}
}

// HandleThing will handle a thing
func (s *Router) HandleThing(thing interface{}) http.HandlerFunc {
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

// NewRoute returns a basic Route interface
func NewRoute(path string, handler http.Handler) Route {
	return &HTTPRoute{
		RoutePath:   path,
		HTTPHandler: handler,
	}
}

// Path gets the path
func (r *HTTPRoute) Path() string {
	return r.RoutePath
}

// Handler gets the handler
func (r *HTTPRoute) Handler() http.Handler {
	return r.HTTPHandler
}
