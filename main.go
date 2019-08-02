package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/pkg/browser"

	"harrybrown.com/app"
	"harrybrown.com/pkg/cmd"
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

var commands = []cmd.Command{
	app.RoutesCmd,
	{
		Syntax:      "time",
		Description: "get the server uptime",
		Run: func() {
			fmt.Println("Server Uptime:", time.Since(serverStart))
		},
	},
	{
		Syntax:      "addr",
		Description: "print out the address that the server is running at",
		Run: func() {
			fmt.Printf("http://%s/\n", addr)
		},
	},
}

var addr string

func main() {
	if debug {
		if *autoOpen {
			addr = open(address, port)
		} else {
			addr = fmt.Sprintf("%s:%s", address, port)
		}

		go cmd.Run(commands)
	} else {
		addr = fmt.Sprintf(":%s", os.Getenv("PORT"))
	}

	if err := server.ListenAndServe(addr); err != nil {
		log.Fatal(err)
	}
}

func open(address, port string) string {
	var addr string
	if address == networkAddr {
		addr = fmt.Sprintf("%s:%s", address, port)
		address = findlocalIP()
	} else {
		addr = fmt.Sprintf("%s:%s", address, port)
	}
	url := fmt.Sprintf("http://%s:%s/", address, port)

	fmt.Println("Running server at", url)
	err := browser.OpenURL(url)
	if err != nil {
		log.Warning(err)
	}
	return addr
}

func findlocalIP() string {
	addrs, _ := net.InterfaceAddrs()
	var ones int
	for _, addr := range addrs {
		if ip, ok := addr.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			ones, _ = ip.Mask.Size()
			if ip.IP.To4() != nil && ones > 16 {
				return ip.IP.String()
			}
		}
	}
	return ""
}
