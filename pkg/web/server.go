package web

import (
	"net/http"

	"harrybrown.com/pkg/log"
)

var (
	// DefaultErrorHandler is the default handler that for errors in the server.
	DefaultErrorHandler http.Handler = http.HandlerFunc(NotFound)

	// DefaultNotFoundHandler is the handler that executes if a page or resourse is not found.
	DefaultNotFoundHandler = http.HandlerFunc(NotFound)

	// DefaultHandlerHook is a hook that alows for the modification of handlers at
	// runtime.
	DefaultHandlerHook func(h http.Handler) http.Handler
)

// ServeMux is an interface that defines compatability with the standard library
// http.ServeMux struct.
type ServeMux interface {
	http.Handler
	Handle(path string, handler http.Handler)
	HandleFunc(path string, handler func(http.ResponseWriter, *http.Request))
}

// Router is an http router.
type Router struct {
	mux         ServeMux
	server      *http.Server
	HandlerHook func(http.Handler) http.Handler
}

// NewRouter creates a new router which contains the default http.ServeMux.
func NewRouter() *Router {
	return &Router{
		mux:         new(http.ServeMux),
		server:      nil,
		HandlerHook: DefaultHandlerHook,
	}
}

// CreateRouter will make a new router from a ServeMux interface.
func CreateRouter(mux ServeMux) *Router {
	return &Router{
		mux:         mux,
		server:      nil,
		HandlerHook: DefaultHandlerHook,
	}
}

// ListenAndServe will run the server.
func (r *Router) ListenAndServe(addr string) error {
	r.server = &http.Server{Addr: addr, Handler: r.mux}
	if r.HandlerHook != nil {
		r.server.Handler = r.HandlerHook(r.mux)
	}

	return r.server.ListenAndServe()
}

// Handle registers the a path and a handler using the standard library interface.
func (r *Router) Handle(path string, h http.Handler) {
	r.mux.Handle(path, h)
}

// HandleFunc will register a new route with a HandlerFunc using the standard library interface.
func (r *Router) HandleFunc(path string, fn http.HandlerFunc) {
	r.mux.Handle(path, http.HandlerFunc(fn))
}

// HandleRoute will handle a route.
func (r *Router) HandleRoute(rt Route) {
	var h http.Handler
	nested, err := rt.Expand()

	if err != nil {
		log.Error(err)

		if e, ok := err.(*ErrorHandler); ok {
			h = e
		} else {
			h = DefaultErrorHandler
		}
	} else {
		h = rt.Handler()
	}
	r.mux.Handle(rt.Path(), h)

	if nested != nil {
		r.HandleRoutes(nested)
	}
}

// HandleRoutes will handle a list of routes.
func (r *Router) HandleRoutes(routes []Route) {
	for _, rt := range routes {
		r.HandleRoute(rt)
	}
}

// AddRoute adds a route to the Router.
func (r *Router) AddRoute(path string, h http.Handler) {
	r.HandleRoute(NewRoute(path, h))
}

// AddRouteFunc adds a route to the Router from a function.
func (r *Router) AddRouteFunc(path string, h http.HandlerFunc) {
	r.HandleRoute(NewRouteFunc(path, h))
}

// HandleThing will handle a thing
//
// Ok, I cant remember why this is here. I think it's just unimplimented
func (r *Router) HandleThing(thing interface{}) http.HandlerFunc {
	// return func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "thing\n")
	// }
	panic("not implimented")
}
