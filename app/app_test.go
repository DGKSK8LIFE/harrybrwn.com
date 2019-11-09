package app

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

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

// {{define "footer" -}}
// <script src="/static/js/home.js"></script>
// {- end}

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

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Error("bad status code")
	}

	page := rr.Body.Bytes()
	if page == nil || len(page) == 0 {
		t.Error("the homepage is empty")
	}
}

func TestRoutes(t *testing.T) {
	for i, route := range Routes {
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
	}
}