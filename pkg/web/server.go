package web

import (
	"net/http"

	"harrybrown.com/pkg/log"
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
	s.mux.Handle(path, http.HandlerFunc(fn))
}

// HandleRoute will handle a route.
func (s *Router) HandleRoute(r Route) {
	var h http.Handler
	nested, err := r.Expand()

	if err != nil {
		log.Error(err)

		if e, ok := err.(*ErrorHandler); ok {
			h = e
		} else {
			h = http.HandlerFunc(NotFound)
		}
	} else {
		h = r.Handler()
	}

	s.mux.Handle(r.Path(), h)

	if nested != nil {
		s.HandleRoutes(nested)
	}
}

// HandleRoutes will handle a list of routes.
func (s *Router) HandleRoutes(routes []Route) {
	for _, r := range routes {
		s.HandleRoute(r)
	}
}

// HandleThing will handle a thing
//
// Ok, I cant remember why this is here. I think it's just unimplimented
func (s *Router) HandleThing(thing interface{}) http.HandlerFunc {
	// return func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "thing\n")
	// }
	panic("not implimented")
}
