package web

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPageTemplatesCollection(t *testing.T) {
	TemplateDir = "/static/templ"
	BaseTemplates = []string{"path/one.html", "path/two.html"}
	p := Page{
		Title:     "Test Page",
		Template:  "path/test/testpage.html",
		RoutePath: "/testpage",
	}
	if p.templateCount() != 3 &&
		p.templateCount() != len(BaseTemplateName)+1 &&
		p.templateCount() != len(p.tmpls()) {
		t.Error("bad template count")
	}

	for _, tmp := range p.tmpls() {
		if !strings.HasPrefix(tmp, TemplateDir) {
			t.Error("templates files given by Page.tmpls() should be in the TemplateDir")
		}
	}

	p.AddTemplateFile("test/dir/file.txt", "test.txt")
	for _, tmp := range p.tmpls() {
		if !strings.HasPrefix(tmp, TemplateDir) {
			t.Error("p.AddTemplateFile should have added the TemplateDir to the front of its input")
		}
	}
}

func newTesters(path string, t *testing.T) (rr *httptest.ResponseRecorder, req *http.Request) {
	var err error
	req, err = http.NewRequest("GET", path, nil)
	if err != nil {
		t.Error("counld not make request")
	}
	rr = httptest.NewRecorder()
	return rr, req
}

func testingPage() *Page {
	return &Page{
		Title:     "Test Page",
		Template:  "",
		RoutePath: "/testing",
		HotReload: false,
	}
}

func TestPageTemplating(t *testing.T) {
	var (
		p        = testingPage()
		err      error
		reqCount = 0
	)
	p.RequestHook = func(self *Page, w http.ResponseWriter, r *http.Request) {
		if len(self.Template) != 0 || self.Title != "Test Page" {
			t.Error("request hook is being given the wrong *Page")
		}
		reqCount++
	}
	err = p.settemplate(BaseTemplateName, "<p>title: {{.Title}}</p>\n<p>path: {{.RoutePath}}</p>")
	if err != nil {
		t.Fatal("template parsing error")
	}

	rr, req := newTesters(p.Path(), t)

	p.Handler().ServeHTTP(rr, req)
	if reqCount != 1 {
		t.Error("request hook is not being executed")
	}

	if rr.Code != 200 {
		t.Error("bad responce code")
	}

	exp := "<p>title: Test Page</p>\n<p>path: /testing</p>"
	if rr.Body.String() != exp {
		t.Error("got:", rr.Body.String(), "expedted:", exp)
	}

	p.Serve = func(w http.ResponseWriter, r *http.Request) {
		reqCount++
	}
	rr, req = newTesters(p.Path(), t)
	p.Handler().ServeHTTP(rr, req)
	if reqCount != 3 {
		t.Error("not all the handlers and hooks in Page are being run upon request")
	}
	p.Serve = nil
}

func TestPageTemplateErrors(t *testing.T) {
	var (
		p    = testingPage()
		err  error
		void = &bytes.Buffer{} // just eat all the bytes
	)

	if _, err = p.WriteTo(void); err == nil {
		t.Error("writing to an io.Writer with no templates should return an error")
	}

	if err = p.settemplate("test", "{{.BadTemplateVar}}"); err != nil {
		t.Error(err)
	}
	if _, err = p.WriteTo(void); err == nil {
		t.Error("writing to a bad template should result in an error")
	}

	rr, req := newTesters(p.Path(), t)
	p.Handler().ServeHTTP(rr, req)
	if rr.Code != 500 {
		t.Error("a request for a page with a bad template should give code 500, got:", rr.Code)
	}
}

func TestPageInitErrors(t *testing.T) {
	p := testingPage()
	err := p.init()
	if err == nil {
		t.Error("init should have returned an error if there are not templates")
	}
	_, err = p.Expand()
	if err == nil {
		t.Error("Expand should also return an error if there are to templates")
	}
}
