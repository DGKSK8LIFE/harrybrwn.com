package web

import (
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

func TestPageTemplating(t *testing.T) {
	var (
		p = &Page{
			Title:     "Test Page",
			Template:  "",
			RoutePath: "/testing",
			HotReload: false,
		}
		err error
	)
	err = p.settemplate(BaseTemplateName, "<p>title: {{.Title}}</p>\n<p>path: {{.RoutePath}}</p>")
	if err != nil {
		t.Fatal("template parsing error")
	}

	req, err := http.NewRequest("GET", p.Path(), nil)
	if err != nil {
		t.Error("failed to make request:", err)
	}
	rr := httptest.NewRecorder()

	p.Handler().ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Error("bad responce code")
	}

	exp := "<p>title: Test Page</p>\n<p>path: /testing</p>"
	if rr.Body.String() != exp {
		t.Error("got:", rr.Body.String(), "expedted:", exp)
	}
}

func TestPageErrors(t *testing.T) {
	var p = &Page{
		Title:     "Page That will fail",
		Template:  "not here",
		RoutePath: "/shouldfail",
	}
	err := p.init()
	if err == nil {
		t.Error("init should have returned an error if there are not templates")
	}
	_, err = p.Expand()
	if err == nil {
		t.Error("Expand should also return an error if there are to templates")
	}
}
