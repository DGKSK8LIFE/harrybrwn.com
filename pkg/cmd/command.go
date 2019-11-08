// Package cmd is an internal packge used for runtime server debugging
package cmd

import (
	"bufio"
	"fmt"
	"os"
)

// Command is a command used in an interface.
type Command struct {
	Syntax      string
	Description string
	Run         func()
}

// DefaultScanner is the package scanner variable.
var DefaultScanner = bufio.NewScanner(os.Stdin)

// Run will run the input interface. Uses the Default Scanner.
func Run(cmds []Command) {
	RunWithScanner(DefaultScanner, cmds)
}

// RunWithScanner runs the input interface using a given scanner.
func RunWithScanner(s *bufio.Scanner, cmds []Command) {
	fmt.Print("> ")
	for s.Scan() {
		if s.Text() == "exit" {
			os.Exit(1)
		}

		if s.Text() == "help" || s.Text() == "h" {
			help(cmds)
		}

		for _, cmd := range cmds {
			if s.Text() == cmd.Syntax {
				cmd.Run()
			}
		}
		fmt.Print("> ")
	}
}

func help(cmds []Command) {
	s := "   "
	fmt.Println("Commands:")
	fmt.Println(s, "exit - exit the program")
	fmt.Println(s, "help - get help")
	for _, cmd := range cmds {
		fmt.Printf("    %s - %s\n", cmd.Syntax, cmd.Description)
	}
	print("\n")
}
