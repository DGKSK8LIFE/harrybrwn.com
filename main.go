package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/pkg/browser"
	"harrybrown.com/app"
	"harrybrown.com/web"
)

var (
	server      = web.NewServer()
	address     = "localhost"
	networkAddr = "0.0.0.0"

	port    = flag.String("port", "8080", "the port to run the server on")
	addrflg = flag.String("address", address, "the address to run the server on")
	network = flag.Bool("network", false, "run the server on the local wifi network (0.0.0.0)")
)

func init() {
	flag.Parse()
	if *network {
		address = networkAddr
	} else if *addrflg != address {
		address = *addrflg
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

func init() {
	web.HandlerHook = app.NewLogger

	server.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	server.HandleRoutes(app.Routes)
}

func main() {
	if err := server.ListenAndServe(open(address, *port)); err != nil {
		log.Fatal(err)
	}
}
