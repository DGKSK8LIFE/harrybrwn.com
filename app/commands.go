package app

import (
	"fmt"
	"time"

	"harrybrown.com/pkg/cmd"
)

var serverStart = time.Now()

// Commands var is a list of debugging commands.
var Commands = []cmd.Command{
	{
		Syntax:      "time",
		Description: "get the server uptime",
		Run: func() {
			fmt.Println("Server Uptime:", time.Since(serverStart))
		},
	},
	{
		Syntax:      "routes",
		Description: "print out all the routes that the server is handling",
		Run: func() {
			for i, r := range Routes {
				fmt.Printf("%d: '%s'\n", i, r.Path())
			}
		},
	},
}
