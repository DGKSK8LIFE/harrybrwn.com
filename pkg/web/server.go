package web

import (
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
	var err error
	for _, r := range routes {
		if err = r.Init(); err != nil {
			s.server.ErrorLog.Printf("Error on Route(\"%s\").Init(): %s\n", r.Path(), err.Error())
		}
		s.Handle(r.Path(), r.Handler())
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
