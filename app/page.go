package app

import (
	"io"
)

// Page is a struct that represents a webpage
type Page struct {
	Title string
	BodyTemplate string
	templates []string
}

// Write will write the webpage to an io.Writer
func (p *Page) Write(w io.Writer) error {
	return execTemplate(
		w,
		newTemplate(p.allTemplates()...),
		// struct { Page *Page } {Page: p},
		p,
	)
}

// AddTemplate will add a template file to the page struct
func (p *Page) AddTemplate(files ...string) {
	p.templates = append(p.templates, files...)
}

func (p *Page) allTemplates() []string {
	return append(p.templates, p.BodyTemplate)
}