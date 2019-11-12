package app

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"harrybrown.com/pkg/log"
	"harrybrown.com/pkg/web"
)

func TestFindUrl(t *testing.T) {
	u, err := url.Parse("http://harrybrwn.com")
	if err != nil {
		t.Error(err)
	}

	var img string
	img = findImage(u, "lg")
	if img != "2250x3000" {
		t.Error("got the wrong image folder from findImage")
	}

	u, err = url.Parse("http://harrybrwn.com/static/img.jpg")
	img = findImage(u, "xs")
	if img != "/static/563x750/img.jpg" {
		t.Error("bad result from findImage:", img)
	}
}

func init() {
	cwd, _ := os.Getwd()
	dir := filepath.Base(cwd)
	if dir == "app" {
		os.Chdir("..")
	}
}

func TestHomePage(t *testing.T) {
	req, err := http.NewRequest("GET", Routes[0].Path(), nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := Routes[0].Handler()

	logger := log.DefaultLogger
	buf := new(bytes.Buffer)
	log.DefaultLogger = log.NewColorLogger(buf, log.Blue)
	handler.ServeHTTP(rr, req) // this is also going to log the error
	if rr.Code != http.StatusInternalServerError {
		t.Error("an uninitialized template should return fail, got:", rr.Code)
	}
	log.DefaultLogger = logger
	if buf.Len() == 0 {
		t.Error("should have logged and error here")
	}

	if route, ok := handler.(*web.Page); ok {
		if inner, err := route.Expand(); inner != nil || err != nil { // needs to read the template files
			t.Error("Pages shouldn't have anything to expand and should succeed")
		}
	}
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Error("bad status code")
	}

	page := rr.Body.Bytes()
	if page == nil || len(page) == 0 {
		t.Error("the homepage is empty")
	}
	if len(page) < 500 {
		t.Error("the home page is too short")
	}
}

func TestAllRoutes(t *testing.T) {
	for i, route := range Routes {
		if route.Path() == "/resume" {
			r := route.(*web.Page)
			r.Data = getResume("./static/data/resume.json")
		}
		if route == nil {
			t.Error("one of your routes get set to 'nil' for some reason")
		}
		if route.Handler() == nil {
			t.Error("route", i, "has no handler, this should be a problem with the 'web' pacage")
		}

		if len(route.Path()) == 0 {
			t.Error("route", i, "has no path")
		}
		if !strings.HasPrefix(route.Path(), "/") {
			t.Error("route", i, "does not begin with a '/'")
		}

		switch r := route.(type) {
		case *web.Page:
			if len(r.Title) == 0 {
				t.Error("route", i, "needs to have a title")
			}
			if len(r.Template) == 0 {
				t.Error("route", i, "needs a template file")
			}
		}

		testGetReq(route, t)
	}
}

func testGetReq(r web.Route, t *testing.T) {
	var (
		err   error
		nodes []web.Route
	)
	if nodes, err = r.Expand(); nodes != nil {
		if err != nil {
			t.Errorf("got error expanding %s: %s\n", r.Path(), err.Error())
		}
		for _, node := range nodes {
			testGetReq(node, t)
		}
	}

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", r.Path(), nil)
	r.Handler().ServeHTTP(rr, req)
	if rr.Code != 200 {
		t.Error("bad response code from", r.Path(), "got", rr.Code)
	}
}

func TestFileServer(t *testing.T) {
	var (
		rr  = httptest.NewRecorder()
		req *http.Request
	)
	fs := web.NewRoute("/static/", NewFileServer("static"))
	req, _ = http.NewRequest("GET", "/static/img/github.svg", nil)

	fs.Handler().ServeHTTP(rr, req)
	if rr.Code != 200 {
		t.Error("bad responce code from file server")
	}
	if len(rr.Body.Bytes()) < 1 {
		t.Error("the file server did not get anything")
	}

	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/static/filenothere.txt", nil)
	fs.Handler().ServeHTTP(rr, req)
	if rr.Code == 200 {
		t.Error("file not found should not give a successful get request")
	}
	if rr.Code != 404 {
		t.Error("why is the code not 404!!")
	}

	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/static/img/me.jpg?size=sm", nil)
	fs.Handler().ServeHTTP(rr, req)
	if rr.Code != 200 {
		t.Error("bad response code")
	}

	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/static/img/me.jpg?size=md", nil)
	fs.Handler().ServeHTTP(rr, req)

	if req.URL.Path != "/static/img/1688x2251/me.jpg" {
		t.Error("wrong url")
	}
	if rr.Code != 200 {
		t.Error("bad response code")
	}
}
