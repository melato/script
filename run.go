package script

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

/** Run 0 or more commands.  If there are more than 1 command, pipe Stdout of each command to Stdin of the next command.
  Any unset Stderr is set to os.Stderr
  Stdin of the first command is set to os.Stdin, unless it is already set
  Stdout of the last command is set to os.Stdout, unless it is already set
  If the global Trace is true, print the commands before running them.
*/
func Run(cmd ...*exec.Cmd) error {
	n := len(cmd)
	if n == 0 {
		return nil
	}
	var closers []io.Closer
	first := cmd[0]
	if first.Stdin == nil {
		first.Stdin = os.Stdin
	}
	last := cmd[n-1]
	if last.Stdout == nil {
		last.Stdout = os.Stdout
	}
	for _, c := range cmd {
		if c.Stderr == nil {
			c.Stderr = os.Stderr
		}
	}
	var err error
	for i, c := range cmd[0 : n-1] {
		next := cmd[i+1]
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
	for i, c := range cmd {
		err = c.Start()
		if err != nil {
			break
		}
		started = i
	}
	if started != len(cmd)-1 {
		for _, c := range cmd[0 : started+1] {
			err := c.Process.Release()
			if err != nil {
				fmt.Println(c.Path, err)
			}
		}
		return err
	}
	waited := -1
	for i, c := range cmd {
		err = c.Wait()
		if err != nil {
			break
		}
		waited = i
	}
	if waited != len(cmd)-1 {
		for _, c := range cmd[waited+1:] {
			err := c.Process.Release()
			if err != nil {
				fmt.Println(c.Path, err)
			}
		}
		return err
	}
	return nil
}
