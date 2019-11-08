package web

import "net/http"

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
