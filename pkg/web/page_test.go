package web

import (
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
