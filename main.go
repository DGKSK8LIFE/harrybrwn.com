package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"harrybrown.com/app"
	"harrybrown.com/pkg/log"
	"harrybrown.com/web"
)

var (
	flags = flag.NewFlagSet("harrybrown.com", flag.ExitOnError)

	server      = web.NewServer()
	address     = "localhost"
	networkAddr = "0.0.0.0"
	port        = "8080"
	debug       = false

	network  = flags.Bool("network", false, "run the server on the local wifi network 0.0.0.0")
	autoOpen = flags.Bool("open", true, "open the webapp in the browser on run")

	serverStart = time.Now()
)

func init() {
	flags.BoolVar(&debug, "debug", debug, "turn on debugging features")
	flags.BoolVar(&debug, "d", debug, "turn on debugging features (shorthand)")
	flags.StringVar(&port, "port", port, "the port to run the server on")
	flags.StringVar(&address, "address", address, "the address to run the server on")
	flags.Parse(os.Args[1:])

	if *network {
		address = networkAddr
	}

	web.HandlerHook = app.NewLogger
	server.HandleRoutes(app.Routes)
}

func main() {
	var addr string
	if debug {
		addr = fmt.Sprintf("%s:%s", address, port)
	} else {
		addr = fmt.Sprintf(":%s", os.Getenv("PORT"))
	}

	if err := server.ListenAndServe(addr); err != nil {
		log.Fatal(err)
	}
}
