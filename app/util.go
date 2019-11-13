package app

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"sort"
	"strings"

	"harrybrown.com/pkg/log"
	"harrybrown.com/pkg/web"
)

// Name is the name of the application.
var Name = "harrybrown.com"

func init() {
	// log.DefaultLogger = log.NewPlainLogger(os.Stdout)
	web.DefaultHandlerHook = NewLogger
	web.DefaultErrorHandler = http.HandlerFunc(NotFound)
}

// NewLogger creates a new logger that will intercept a handler and replace it
// with one that has logging functionality.
func NewLogger(h http.Handler) http.Handler {
	return &pageLogger{
		wrapped: h,
		l:       log.NewPlainLogger(os.Stdout),
	}
}

type pageLogger struct {
	wrapped http.Handler
	l       log.PrintLogger
}

func (p *pageLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.l.Printf("%s %s%s\n", r.Method, r.Host, r.URL)
	p.wrapped.ServeHTTP(w, r)
}

// NotFound handles requests that generate a 404 error
func NotFound(w http.ResponseWriter, r *http.Request) {
	var tmplNotFound = template.Must(template.ParseFiles(
		"templates/pages/404.html",
		"templates/index.html",
	))
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if err := tmplNotFound.ExecuteTemplate(w, "base", nil); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

// ParseFlags parses flags.
func ParseFlags() {
	defer RecoverFlagHelpErr()
	flag.Parse()
}

// StringFlag adds a string flag with a shorthand.
func StringFlag(ptr *string, name, desc string) {
	flag.StringVar(ptr, name, *ptr, desc)
	flag.StringVar(ptr, name[:1], *ptr, desc)
}

// BoolFlag adds a boolean flag with a shorthand.
func BoolFlag(ptr *bool, name, desc string) {
	flag.BoolVar(ptr, name, *ptr, desc)
	flag.BoolVar(ptr, name[:1], *ptr, desc)
}

// RecoverFlagHelpErr will gracfully end the program if the help flag is given.
func RecoverFlagHelpErr() {
	if r := recover(); r != nil {
		err, ok := r.(error)
		if !ok {
			return
		}
		if err == flag.ErrHelp {
			os.Exit(0)
		} else {
			fmt.Println("Error:", err.Error())
			os.Exit(2)
		}
	}
}

type appFlag struct {
	flag      *flag.Flag
	name      string
	shorthand string
	usage     string
	fmtlen    int
}

func (af *appFlag) len() int {
	return len(af.name)
}

func (af *appFlag) useline() string {
	if len(af.shorthand) != 0 {
		return fmt.Sprintf("-%s, -%s", af.shorthand, af.name)
	}
	return fmt.Sprintf("    -%s", af.name)
}

func (af *appFlag) defval() string {
	deflt := ""
	if af.flag != nil {
		deflt = fmt.Sprintf("(default: %s)", af.flag.DefValue)
	}
	return deflt
}

func (af *appFlag) String() string {
	return fmt.Sprintf("  %s%s%s %s",
		af.useline(),
		strings.Repeat(" ", 4+af.fmtlen-af.len()),
		af.usage,
		af.defval())
}

func newflag(flg *flag.Flag, shorthand string) *appFlag {
	return &appFlag{
		flag:      flg,
		name:      flg.Name,
		usage:     flg.Usage,
		shorthand: shorthand,
	}
}

func getFlags() ([]*appFlag, int) {
	fmap := make(map[string]*appFlag)
	flag.VisitAll(func(flg *flag.Flag) {
		if _, inmap := fmap[flg.Usage]; !inmap {
			fmap[flg.Usage] = new(appFlag)
			if len(flg.Name) == 1 {
				fmap[flg.Usage].shorthand = flg.Name
			}
		}
		fmap[flg.Usage] = newflag(flg, fmap[flg.Usage].shorthand)
	})
	length := len(fmap)

	flgs := make([]*appFlag, length)
	i := 0
	for _, fl := range fmap {
		flgs[i] = fl
		i++
	}

	max := flgs[0].len()
	for i = 1; i < length; i++ {
		if flgs[i].len() > max {
			max = flgs[i].len()
		}
	}

	sort.Slice(flgs, func(i, j int) bool {
		return flgs[i].name[0] < flgs[j].name[0]
	})
	return flgs, max
}

var helpFlag = &appFlag{name: "help", shorthand: "h", usage: "get help", fmtlen: 4}

func init() {
	out := flag.CommandLine.Output()

	flag.CommandLine.Usage = func() {
		flags, maxlen := getFlags()
		helpFlag.fmtlen = maxlen
		flags = append(flags, helpFlag)

		fmt.Fprintf(out, "Usage of %s:\n", Name)
		for _, v := range flags {
			v.fmtlen = maxlen
			fmt.Fprintln(out, v)
		}
		fmt.Fprint(out, "\n")
	}
}
