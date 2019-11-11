package web

import (
	"encoding/json"
	"net/http"

	"harrybrown.com/pkg/log"
)

// JSONRoute is an api route that returns json.
type JSONRoute struct {
	APIPath string
	Run     func(http.ResponseWriter, *http.Request) interface{}
}

// NewJSONRoute creates a new json route.
func NewJSONRoute(path string, fn func(http.ResponseWriter, *http.Request) interface{}) *JSONRoute {
	return &JSONRoute{
		APIPath: path,
		Run:     fn,
	}
}

// Path will return the route path.
func (j *JSONRoute) Path() string {
	return j.APIPath
}

func (j *JSONRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	e := json.NewEncoder(w)
	if err := e.Encode(j.Run(w, r)); err != nil {
		log.Error("Json Error:", err.Error())
	}
}

// Handler will return the handler.
func (j *JSONRoute) Handler() http.Handler {
	return j
}

// Expand does nothing for json routes.
func (j *JSONRoute) Expand() []Route {
	return nil
}

var (
	_ Route        = (*JSONRoute)(nil)
	_ http.Handler = (*JSONRoute)(nil)
)
