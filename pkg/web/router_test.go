package web

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var hello = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello tests")
})

var testingRoutes = []Route{
	NewRouteFunc("/hello", hello),
	NewRoute("/test", http.NotFoundHandler()),
}

func TestRouter(t *testing.T) {
	r := NewRouter()
	r.HandleRoutes(testingRoutes)

	ts := httptest.NewServer(r.mux)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/hello")
	if err != nil {
		t.Error(err)
	}
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	if string(response) != "hello tests" {
		t.Error("got the wrong response:", string(response))
	}

	res, err = http.Get(ts.URL + "/test")
	if res.StatusCode != 404 {
		t.Error("Wrong status code")
	}
}

func TestNestedRoute(t *testing.T) {
	r := NewRouter()
	tests := NewNestedRoute("/testing", testingRoutes...)
	r.HandleRoute(tests)
	ts := httptest.NewServer(r.mux)

	res, err := http.Get(ts.URL + "/testing/hello")
	if err != nil {
		t.Error(err)
	}
	text, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	if string(text) != "hello tests" {
		t.Error("get wrong response text:", string(text))
	}

	res, err = http.Get(ts.URL + "/testing/nothere")
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 404 {
		t.Error("wrong status code:", res.StatusCode)
	}
}
