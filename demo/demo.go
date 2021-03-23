package main

import (
	"fmt"

	"melato.org/script"
)

func space() {
	fmt.Println()
	fmt.Println()
}

func main() {
	s := script.Script{Trace: true}

	fmt.Println("Run:")
	s.Run("ls", "-s")

	space()
	fmt.Println("redirect to file:")
	s.Cmd("ls").ToFile("ls.out")

	space()
	fmt.Println("input string:")
	s.Cmd("cat").InputString("hello").Run()

	space()
	fmt.Println("input bytes:")
	s.Cmd("cat").InputBytes([]byte("bytes")).Run()

	space()
	fmt.Println("input file:")
	s.Cmd("cat").InputFile("ls.out").Run()

	space()
	pwd := s.Cmd("pwd").ToString()
	fmt.Println("output to string: '" + pwd + "'")

	space()
	fmt.Println("pipe two commands:")
	s.Cmd("ls", "-1").PipeTo("sort", "-r").Run()

	if s.HasError() {
		fmt.Println(s.Error())
	}
}
