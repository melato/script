package main

import (
	"os/exec"

	"melato.org/script"
)

func main() {
	commands := []*exec.Cmd{
		exec.Command("ls", "-1"),
		exec.Command("sort", "-r")}
	script.PrintPipeline(commands...)
	script.Pipeline(commands...)
}
