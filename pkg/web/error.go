package web

import (
	"fmt"
	"html/template"
	"net/http"
	"runtime"

	"harrybrown.com/pkg/log"
)

// ErrorHandler is an error type for internal website errors.
type ErrorHandler struct {
	msg      string
	status   int
	file     string
	funcname string
	line     int
}

// NotFound handles requests that generate a 404 error
func NotFound(w http.ResponseWriter, r *http.Request) {
	var tmplNotFound = template.Must(template.ParseFiles(
		getfile("pages/404.html"),
		getfile("index.html"),
	))
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if err := tmplNotFound.ExecuteTemplate(w, "base", nil); err != nil {
		fmt.Fprintf(w, "%s", err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

// Error creates a new error.
func Error(status int, msg string) error {
	pc, file, line, _ := runtime.Caller(1)
	e := &ErrorHandler{
		msg:      msg,
		status:   status,
		file:     file,
		line:     line,
		funcname: runtime.FuncForPC(pc).Name(),
	}
	return e
}

// Errorf create an error with a formatted message.
func Errorf(status int, format string, vars ...interface{}) error {
	pc, file, line, _ := runtime.Caller(1)

	e := &ErrorHandler{
		msg:      fmt.Sprintf(format, vars...),
		status:   status,
		file:     file,
		line:     line,
		funcname: runtime.FuncForPC(pc).Name(),
	}
	return e
}

func (h *ErrorHandler) Error() string {
	return fmt.Sprintf("(%s:%d %s()) %s\n", h.file, h.line, h.funcname, h.msg)
}

func (h *ErrorHandler) log() {
	log.Error(h.Error())
}

func (h *ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ServeError(w, h.status)
}

var (
	_ error        = (*ErrorHandler)(nil)
	_ http.Handler = (*ErrorHandler)(nil)
)

// ServeError serves a generic http error page.
func ServeError(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	t, err := template.New("").Parse(`
{{define "header" -}}
<title>{{.Title}}</title>
<style>h2, .ErrorMsg { text-align: center; }</style>
{{- end}}
{{define "body" -}}
<div class="container">
	<h2>{{.Status}} Something Went Wrong</h2>
	<div class="ErrorMsg">
		<p>Sorry, I must have broken something.</p>
		<p hidden>if you can see this then i am sorry</p>
	</div>
</div>
{{- end}}
{{define "footer" -}}{{- end}}`)
	if err != nil {
		log.Println("Error when serving error page:", err)
	}
	t, err = t.ParseFiles(getfile("index.html"))
	if err != nil {
		log.Println("Error when serving error page:", err)
	}
	if err = t.ExecuteTemplate(w, "base", struct {
		Title  string
		Status int
	}{
		Title:  "Error",
		Status: status,
	}); err != nil {
		log.Println("Error when serving error page:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
