package main

import (
	"flag"
	"fmt"
	"os"

	"harrybrown.com/app"
	"harrybrown.com/pkg/cmd"
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

	router.HandleRoutes(app.Routes)
}

func main() {
	if !debug {
		port = os.Getenv("PORT") // use heroku's port variable
	} else {
		log.Printf("running on %s:%s\n", address, port)

		app.Commands = append(app.Commands, cmd.Command{
			Syntax:      "addr",
			Description: "show server address",
			Run:         func() { fmt.Printf("%s:%s\n", address, port) },
		})

		// scans stdin and runs the commands given alongside the server
		go cmd.Run(app.Commands)
	}

	if err := router.ListenAndServe(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal(err)
	}
}
