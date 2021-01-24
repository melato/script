package script

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func Println(cmd *exec.Cmd, extraArgs ...string) {
	cmdLen := len(cmd.Args)
	args := make([]interface{}, cmdLen+len(extraArgs))
	args[0] = cmd.Path
	for i, arg := range cmd.Args {
		args[i] = arg
	}
	for i, arg := range extraArgs {
		args[cmdLen+i] = arg
	}
	fmt.Println(args...)
}

func PrintPipeline(commands ...*exec.Cmd) {
	n := len(commands)
	if n == 0 {
		return
	}
	for _, c := range commands[0 : n-1] {
		Println(c, "|")
	}
	Println(commands[n-1])
}

func Pipeline(commands ...*exec.Cmd) error {
	var closers []io.Closer

	n := len(commands)
	if n == 0 {
		return nil
	}
	first := commands[0]
	if first.Stdin == nil {
		first.Stdin = os.Stdin
	}
	last := commands[n-1]
	if last.Stdout == nil {
		last.Stdout = os.Stdout
	}
	for _, c := range commands {
		if c.Stderr == nil {
			c.Stderr = os.Stderr
		}
	}
	var err error
	for i, c := range commands[0 : n-1] {
		next := commands[i+1]
		var in io.ReadCloser
		in, err = c.StdoutPipe()
		if err != nil {
			break
		}
		next.Stdin = in
		closers = append(closers, in)
	}
	if err != nil {
		for _, closer := range closers {
			closer.Close()
		}
		return err
	}
	started := -1
	for i, c := range commands {
		err = c.Start()
		if err != nil {
			break
		}
		started = i
	}
	if started != len(commands)-1 {
		for _, c := range commands[0 : started+1] {
			err := c.Process.Release()
			if err != nil {
				fmt.Println(c.Path, err)
			}
		}
		return err
	}
	waited := -1
	for i, c := range commands {
		err = c.Wait()
		if err != nil {
			break
		}
		waited = i
	}
	if waited != len(commands)-1 {
		for _, c := range commands[waited+1:] {
			err := c.Process.Release()
			if err != nil {
				fmt.Println(c.Path, err)
			}
		}
		return err
	}
	return nil

}
