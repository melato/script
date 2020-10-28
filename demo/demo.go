package main

import (
	"fmt"

	"melato.org/script"
)

func main() {
	s := script.Script{Trace: true}

	s.Run("ls")

	s.Cmd("ls").ToFile("ls.out")
	pwd := s.Cmd("pwd").ToString()
	fmt.Println("pwd: '" + pwd + "'")

	s.Cmd("ls", "-1").Pipe(s.Cmd("sort", "-r")).Run()

	if s.Error != nil {
		fmt.Println(s.Error)
	}
}
