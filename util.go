package script

import (
	"bufio"
	"bytes"
	"os/exec"
)

func BytesToLines(out []byte) []string {
	scanner := bufio.NewScanner(bytes.NewReader(out))
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

func Pipe(cmds ...*exec.Cmd) error {
	n := len(cmds)
	var err error
	for i := 0; i < n; i++ {
		c := cmds[i]
		if i < n-1 {
			c2 := cmds[1]
			c2.Stdin, err = c.StdoutPipe()
			if err != nil {
				return err
			}
		}
	}
	for _, c := range cmds {
		if err := c.Start(); err != nil {
			return err
		}

	}
	for _, c := range cmds {
		if err := c.Wait(); err != nil {
			return err
		}

	}
	return nil
}
