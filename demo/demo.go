package main

import (
	"fmt"

	"melato.org/script"
)

func main() {
	var s script.Script
	s.Trace = true

	s.Cmd("ls").ToFile("ls.out")
	pwd := s.Cmd("pwd").ToString()
	fmt.Println("pwd: '" + pwd + "'")

	s.Cmd("ls").Pipe(s.Cmd("sort", "-r")).Run()

	if s.Error != nil {
		fmt.Println(s.Error)
	}
}
