package web

import "net/http"

// Route defines an interface for various routes to be used in the Server struct.
type Route interface {
	// Path should return the path for that specific route.
	Path() string

	// Handler returns the handler that will be used to server up http responces.
	Handler() http.Handler

	// Init should initialize the Route. In relation to the web.Router, this method
	// is called for every Route given to Router.HandleRoutes.
	Init() error
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

// Init will initialize the HTTPRoute.
//
// (this function currently does nothing)
func (r *HTTPRoute) Init() error {
	return nil
}
