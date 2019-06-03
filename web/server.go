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

// HandleThing will handle a thing
func (s *Server) HandleThing(thing interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "thing\n")
	}
}
