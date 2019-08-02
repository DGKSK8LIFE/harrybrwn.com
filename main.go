package main

import (
	"flag"
	"fmt"
	"os"

	"harrybrown.com/app"
	"harrybrown.com/pkg/log"
	"harrybrown.com/pkg/web"
)

var (
	flags = flag.NewFlagSet("harrybrown.com", flag.ExitOnError)

	router      = web.NewRouter()
	address     = "localhost"
	networkAddr = "0.0.0.0"
	port        = "8080"
	debug       = false
)

func init() {
	flags.BoolVar(&debug, "debug", debug, "turn on debugging features")
	flags.BoolVar(&debug, "d", debug, "turn on debugging features (shorthand)")
	flags.StringVar(&port, "port", port, "the port to run the server on")
	flags.StringVar(&address, "address", address, "the address to run the server on")
	flags.Parse(os.Args[1:])

	web.HandlerHook = app.NewLogger
	router.HandleRoutes(app.Routes)
}

func main() {
	var addr string
	if debug {
		addr = fmt.Sprintf("%s:%s", address, port)
	} else {
		addr = fmt.Sprintf(":%s", os.Getenv("PORT"))
	}

	if err := router.ListenAndServe(addr); err != nil {
		log.Fatal(err)
	}
}
