package main

import (
	"fmt"

	"melato.org/script"
)

func main() {
	var s script.Script

	s.Cmd("ls").ToFile("ls.out")

	//s.Cmd("ls").Pipe(s.Cmd("sort")).Run()

	pwd := s.Cmd("pwd").ToString()
	fmt.Println("pwd: '" + pwd + "'")

	if s.Error != nil {
		fmt.Println(s.Error)
	}
}
