package script

/**
scripting utilities:
- Run programs
- Run pipelines
- Redirect output to file
- Copy output to a string or []string

- Maintains an internal error state, so the user does not have to check for an error after running each command.
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
	Trace bool  // print exec arguments to stdout
	Error error // The first error encountered.  If not nil, do not execute anything else.
}

type Cmd struct {
	Cmd         *exec.Cmd
	inputCmd    *Cmd
	IsRun       bool
	mergeStderr bool
	script      *Script
}

/** Run a command or pipeline */
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
		if t.script.InError(err) {
			return
		}
	}
	for _, c := range commands {
		if t.script.InError(c.Cmd.Start()) {
			return
		}
	}
	for _, c := range commands {
		if t.script.InError(c.Cmd.Wait()) {
			return
		}
	}
	return
}

/** Run and redirect output to a Writer */
func (t *Cmd) ToWriter(out io.Writer) {
	if t.script.HasError() {
		return
	}
	t.Cmd.Stdout = out
	t.Run()
}

/** Run and redirect output to a file */
func (t *Cmd) ToFile(file string) {
	if t.script.HasError() {
		return
	}
	f, err := os.Create(file)
	if t.script.InError(err) {
		return
	}
	defer f.Close()
	t.ToWriter(f)
}

/** Run and return the output */
func (t *Cmd) ToBytes() []byte {
	var buf bytes.Buffer
	t.ToWriter(&buf)
	return buf.Bytes()
}

/** Run and return the output as a string */
func (t *Cmd) ToString() string {
	return strings.TrimSpace(string(t.ToBytes()))
}

/** Run and return the output as a []string */
func (t *Cmd) ToLines() []string {
	return BytesToLines(t.ToBytes())
}

func (t *Cmd) Pipe(to *Cmd) *Cmd {
	to.inputCmd = t
	return to
}

/** Return true if an error has happened */
func (t *Script) HasError() bool {
	return t.Error != nil
}

/** Check if the argument is an error, or whether the script already has an error.
  If the argument is an error and it is the first error encountered, it becomes the script's error.
  Return true if the script has an error, either the given error or a previous one.  */
func (t *Script) InError(err error) bool {
	if t.HasError() {
		return true
	}
	if err != nil {
		t.Error = err
		return true
	}
	return false
}

/** Create a command without running it.  The command can be executed or piped to another command. */
func (t *Script) Cmd(name string, args ...string) *Cmd {
	r := &Cmd{}
	r.script = t
	if t.Error != nil {
		return r
	}
	path, err := exec.LookPath(name)
	if t.InError(err) {
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

/** Create a command and run it */
func (t *Script) Run(name string, args ...string) {
	t.Cmd(name, args...).Run()
}

func (t *Cmd) PipeTo(name string, args ...string) *Cmd {
	return t.Pipe(t.script.Cmd(name, args...))
}

func (t *Cmd) MergeStderr() {
	t.mergeStderr = true
}

func (t *Cmd) Script() *Script {
	return t.script
}
