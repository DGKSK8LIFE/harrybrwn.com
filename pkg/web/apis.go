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
	Static  func() interface{}
}

// APIRoute creates a new json route that has access to the response and request.
func APIRoute(path string, fn func(http.ResponseWriter, *http.Request) interface{}) *JSONRoute {
	return &JSONRoute{
		APIPath: path,
		Run:     fn,
	}
}

// StaticAPIRoute creates a new json route the has no access to the response and request.
func StaticAPIRoute(path string, fn func() interface{}) *JSONRoute {
	return &JSONRoute{
		APIPath: path,
		Static:  fn,
	}
}

// ServeJSON will turn an interface into json and write it to a response-writer.
func ServeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Error("Json Error:", err.Error())
		ServeError(w, 500)
	}
}

// Path will return the route path.
func (j *JSONRoute) Path() string {
	return j.APIPath
}

func (j *JSONRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var data interface{}

	if j.Static != nil {
		data = j.Static()
	} else if j.Run != nil {
		data = j.Run(w, r)
	}

	ServeJSON(w, data)
}

// Handler will return the handler.
func (j *JSONRoute) Handler() http.Handler {
	return j
}

// Expand does nothing for json routes.
func (j *JSONRoute) Expand() ([]Route, error) {
	return nil, nil
}

var (
	_ Route        = (*JSONRoute)(nil)
	_ http.Handler = (*JSONRoute)(nil)
)
