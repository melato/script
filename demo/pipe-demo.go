package main

import (
	"fmt"
	"os"
	"os/exec"

	"melato.org/script"
)

func main() {
	commands := []*exec.Cmd{
		exec.Command("ls", "-1"),
		exec.Command("sort", "-r")}
	if len(os.Args) > 1 {
		commands = append(commands, exec.Command(os.Args[1], os.Args[2:]...))
	}

	script.Trace = true
	err := script.Run(commands...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
