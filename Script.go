package script

/**
scripting utilities
Run programs
Run pipelines
Redirect output to file
Copy output to memory
Split output lines
*/

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Script struct {
	Trace bool `usage:"print call arguments to stdout"`
	Error error
}

type Cmd struct {
	Cmd         *exec.Cmd
	inputCmd    *Cmd
	IsRun       bool
	mergeStderr bool
	script      *Script
}

func (t *Cmd) Run() {
	if t.IsRun {
		return
	}
	if t.script.HasError() {
		return
	}
	if t.Cmd.Stdout == nil {
		t.Cmd.Stdout = os.Stdout
	}
	if t.mergeStderr {
		t.Cmd.Stderr = t.Cmd.Stdout
	} else {
		t.Cmd.Stderr = t.Cmd.Stdout
	}
	var commands []*Cmd
	for c := t; c != nil; c = c.inputCmd {
		commands = append(commands, c)
	}
	for i, j := 0, len(commands)-1; i < j; i, j = i+1, j-1 {
		commands[i], commands[j] = commands[j], commands[i]
	}
	commands[0].Cmd.Stdin = os.Stdin
	n := len(commands)
	for i := 0; i < n-1; i++ {
		var err error
		commands[i+1].Cmd.Stdin, err = commands[i].Cmd.StdoutPipe()
		if t.script.checkError(err) {
			return
		}
	}
	for _, c := range commands {
		if t.script.checkError(c.Cmd.Start()) {
			return
		}
	}
	for _, c := range commands {
		if t.script.checkError(c.Cmd.Wait()) {
			return
		}
	}
	return
}

func (t *Cmd) ToWriter(out io.Writer) {
	if t.script.HasError() {
		return
	}
	t.Cmd.Stdout = out
	t.Run()
}

func (t *Cmd) ToFile(file string) {
	if t.script.HasError() {
		return
	}
	f, err := os.Create(file)
	if t.script.checkError(err) {
		return
	}
	defer f.Close()
	t.ToWriter(f)
}

func (t *Cmd) ToBytes() []byte {
	var buf bytes.Buffer
	t.ToWriter(&buf)
	return buf.Bytes()
}

func (t *Cmd) ToString() string {
	return strings.TrimSpace(string(t.ToBytes()))
}

func (t *Cmd) ToLines() []string {
	return BytesToLines(t.ToBytes())
}

func (t *Cmd) Pipe(to *Cmd) *Cmd {
	to.inputCmd = t
	return to
}

func (t *Script) HasError() bool {
	return t.Error != nil
}

func (t *Script) checkError(err error) bool {
	if t.HasError() {
		return true
	}
	if err != nil {
		t.Error = err
		return true
	}
	return false
}

func (t *Script) Cmd(name string, args ...string) *Cmd {
	r := &Cmd{}
	r.script = t
	if t.Error != nil {
		return r
	}
	path, err := exec.LookPath(name)
	if t.checkError(err) {
		return r
	}
	cargs := []string{name}
	cargs = append(cargs, args...)
	if t.Trace {
		fmt.Println(path, cargs)
	}
	r.Cmd = &exec.Cmd{Path: path, Args: cargs}
	return r
}

func (t *Cmd) MergeStderr() {
	t.mergeStderr = true
}
