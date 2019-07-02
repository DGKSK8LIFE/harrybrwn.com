package app

import (
	"fmt"

	"harrybrown.com/pkg/cmd"
)

// RoutesCmd is the command that prints out the roues.
var RoutesCmd = cmd.Command{
	Syntax:      "routes",
	Description: "print out all the routes that the server is handling",
	Run: func() {
		for i, r := range Routes {
			fmt.Printf("%d: '%s'\n", i, r.Path())
		}
	},
}
