package web

import (
	"net/http"
	"path/filepath"
)

// Route defines an interface for various routes to be used in the Server struct.
type Route interface {
	// Path should return the path for that specific route.
	Path() string

	// Handler returns the handler that will be used to server up http responces.
	Handler() http.Handler

	// Expand will expand the route if it has any nested routes. Returns nil if there
	// are no nested routes.
	Expand() ([]Route, error)
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

// NewRouteFunc creates a route from an http.HandlerFunc.
func NewRouteFunc(path string, handler http.HandlerFunc) Route {
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

// Expand returuns nothing for a baseic HTTPRoute.
func (r *HTTPRoute) Expand() ([]Route, error) {
	return nil, nil
}

// NestedRoute is a Route that can serve its own requests but also holds.
type NestedRoute struct {
	BasePath    string
	BaseHandler http.Handler
	routes      []Route
}

// NewNestedRoute creates a new Nested route
func NewNestedRoute(path string, routes ...Route) *NestedRoute {
	rts := make([]Route, len(routes))
	for i, r := range routes {
		rts[i] = &innerRoute{basepath: path, inner: r}
	}
	return &NestedRoute{
		BasePath:    path,
		BaseHandler: http.HandlerFunc(NotImplimented),
		routes:      rts,
	}
}

// SetHandler sets the BaseHander field of the NestedRoute and returns
// the NestedRoute.
func (nr *NestedRoute) SetHandler(h http.Handler) *NestedRoute {
	nr.BaseHandler = h
	return nr
}

// SetHandlerFunc sets the BaseHandler field as the handlerfunc given.
func (nr *NestedRoute) SetHandlerFunc(fn http.HandlerFunc) *NestedRoute {
	nr.BaseHandler = fn
	return nr
}

// Path returns the base path for the nested list of Routes.
func (nr *NestedRoute) Path() string {
	return nr.BasePath
}

// Handler returns the handler that server that base path.
func (nr *NestedRoute) Handler() http.Handler {
	return nr.BaseHandler
}

// Expand will return the list of nested Routes
func (nr *NestedRoute) Expand() ([]Route, error) {
	return nr.routes, nil
}

// AddRoute will add a nested Route.
func (nr *NestedRoute) AddRoute(r Route) *NestedRoute {
	nr.routes = append(nr.routes, &innerRoute{basepath: nr.BasePath, inner: r})
	return nr
}

var (
	_ Route = (*NestedRoute)(nil)
	_ Route = (*innerRoute)(nil)
)

type innerRoute struct {
	basepath string
	inner    Route
}

func (ir *innerRoute) Path() string {
	return filepath.Join(ir.basepath, ir.inner.Path())
}

func (ir *innerRoute) Handler() http.Handler {
	return ir.inner.Handler()
}

func (ir *innerRoute) Expand() ([]Route, error) {
	return ir.inner.Expand()
}
