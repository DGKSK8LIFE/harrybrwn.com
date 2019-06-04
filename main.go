package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/browser"
	"harrybrown.com/app"
	"harrybrown.com/web"
)

const (
	port    = "8080"
	address = "localhost"
)

var (
	server = web.NewServer()
	addr   string
)

func init() {
	web.HandlerHook = app.NewLogger

	server.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	server.HandleRoutes(app.Routes)
}

func main() {
	addr := fmt.Sprintf("%s:%s", address, port)
	url := fmt.Sprintf("http://%s/", addr)

	if err := browser.OpenURL(url); err != nil {
		log.Println(err)
	}

	if err := server.ListenAndServe(addr); err != nil {
		log.Fatal(err)
	}
}
