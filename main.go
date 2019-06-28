package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/pkg/browser"
	"harrybrown.com/app"
	"harrybrown.com/web"
)

var (
	flags = flag.NewFlagSet("harrybrown.com", flag.ExitOnError)

	server      = web.NewServer()
	address     = "localhost"
	networkAddr = "0.0.0.0"
	port        = "8080"

	network  = flags.Bool("network", false, "run the server on the local wifi network (0.0.0.0)")
	autoOpen = flags.Bool("open", true, "open the webapp in the browser on run")
)

func init() {
	flags.StringVar(&port, "port", port, "the port to run the server on")
	flags.StringVar(&address, "address", address, "the address to run the server on")

	if len(os.Args) > 2 {
		if os.Args[1] == "-h" || os.Args[1] == "-help" {
			flags.Usage()
			os.Exit(0)
		}
	}

	flags.Parse(os.Args[1:])
	if *network {
		address = networkAddr
	}

	web.HandlerHook = app.NewLogger

	server.Handle("/static/", app.NewFileServer("static"))
	server.HandleRoutes(app.Routes)
}

func main() {
	var addr string
	if *autoOpen {
		addr = open(address, port)
	} else {
		addr = fmt.Sprintf("%s:%s", address, port)
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
	err := browser.OpenURL(fmt.Sprintf("http://%s:%s/", address, port))
	if err != nil {
		log.Fatal(err)
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
