package web

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"harrybrown.com/pkg/log"
)

var hello = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello tests")
})

var testingRoutes = []Route{
	NewRouteFunc("/hello", hello),
	NewRoute("/test", http.NotFoundHandler()),
}

func errif(t *testing.T, e error) {
	if e != nil {
		t.Error(e.Error())
	}
}

func TestRouter(t *testing.T) {
	r := NewRouter()
	r.HandleRoutes(testingRoutes)
	r.Handle("/", http.NotFoundHandler())

	ts := httptest.NewServer(r.mux)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/hello")
	errif(t, err)
	response, err := ioutil.ReadAll(res.Body)
	errif(t, err)
	if string(response) != "hello tests" {
		t.Error("got the wrong response:", string(response))
	}

	res, err = http.Get(ts.URL + "/test")
	if res.StatusCode != 404 {
		t.Error("Wrong status code")
	}
	res, err = http.Get(ts.URL + "/")
	if res.StatusCode != 404 {
		t.Error("Wrong status code")
	}
}

func TestHandleFunc(t *testing.T) {
	r := NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		fmt.Fprint(w, "this is a test")
	})
	ts := httptest.NewServer(r.mux)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/")
	errif(t, err)
	if res.StatusCode != 201 {
		t.Error("bad status code, got:", res.StatusCode)
	}
	text, err := ioutil.ReadAll(res.Body)
	errif(t, err)
	if string(text) != "this is a test" {
		t.Error("bad output")
	}
}

func TestNestedRoute(t *testing.T) {
	r := NewRouter()
	tests := NewNestedRoute("/testing", testingRoutes...)
	r.HandleRoute(tests)
	ts := httptest.NewServer(r.mux)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/testing/hello")
	errif(t, err)

	text, err := ioutil.ReadAll(res.Body)
	errif(t, err)
	if string(text) != "hello tests" {
		t.Error("get wrong response text:", string(text))
	}

	res, err = http.Get(ts.URL + "/testing/nothere")
	errif(t, err)
	if res.StatusCode != 404 {
		t.Error("wrong status code:", res.StatusCode)
	}
	text, err = ioutil.ReadAll(res.Body)
	errif(t, err)
}

func TestJSONRoute(t *testing.T) {
	r := NewRouter()
	r.HandleRoute(
		APIRoute("/test", func(w http.ResponseWriter, r *http.Request) interface{} {
			return struct {
				One string
				Two int
			}{One: "one", Two: 2}
		}),
	)
	r.HandleRoute(StaticAPIRoute("/test/static", func() interface{} {
		return map[string]interface{}{"One": "one", "Two": 2}
	}))

	ts := httptest.NewServer(r.mux)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/test")
	errif(t, err)
	if res.StatusCode != 200 {
		t.Error("bad status code")
	}
	text, err := ioutil.ReadAll(res.Body)
	errif(t, err)
	if string(text) != "{\"One\":\"one\",\"Two\":2}\n" {
		t.Error("Got unexpected output:", string(text))
	}
	res = nil
	text = nil

	res, err = http.Get(ts.URL + "/test/static")
	errif(t, err)
	if res.StatusCode != 200 {
		t.Error("bad status code")
	}
	text, err = ioutil.ReadAll(res.Body)
	errif(t, err)
	if string(text) != "{\"One\":\"one\",\"Two\":2}\n" {
		t.Error("Got unexpected output:", string(text))
	}
}

type badroute struct {
	p string
	e error
}

func (b *badroute) Path() string {
	return b.p
}

func (b *badroute) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprint(w, "hello?")
	})
}

func (b *badroute) Expand() ([]Route, error) {
	return nil, b.e
}

var _ Route = (*badroute)(nil)

func defaultErrorHandlerOutput() string {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "", nil)
	DefaultErrorHandler.ServeHTTP(rr, req)
	return rr.Body.String()
}

func TestExpandError(t *testing.T) {
	logger := log.DefaultLogger
	buf := new(bytes.Buffer)
	log.DefaultLogger = log.NewColorLogger(buf, log.Blue)
	r := CreateRouter(new(http.ServeMux))
	r.HandleRoute(&badroute{p: "/", e: errors.New("badroute expand error")})
	r.HandleRoute(&badroute{p: "/error", e: Error(501, "should fail")})
	if buf.Len() == 0 {
		t.Error("should have printed an error")
	}
	if !strings.HasPrefix(buf.String(), string(log.Red)) {
		t.Error("the output should have been red")
	}

	ts := httptest.NewServer(r.mux)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/")
	errif(t, err)
	if res.StatusCode == 200 {
		t.Error("the badroute should not have run its handler successfully")
	}
	if res.StatusCode != 404 {
		t.Error("status should be 404")
	}
	text, err := ioutil.ReadAll(res.Body)
	errif(t, err)
	if string(text) != defaultErrorHandlerOutput() {
		t.Error("got the wrong error handleing output")
	}
	res = nil
	text = nil

	res, err = http.Get(ts.URL + "/error")
	errif(t, err)
	if res.StatusCode != 501 {
		t.Error("should have gotten 501 as the status code got:", res.StatusCode)
	}
	text, err = ioutil.ReadAll(res.Body)
	errif(t, err)
	if !bytes.Contains(text, []byte("Response Code 501")) {
		t.Error("the output page should contain error response code")
		fmt.Println(string(text))
	}
	log.DefaultLogger = logger
}
