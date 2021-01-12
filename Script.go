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
	Trace  bool // print exec arguments to stdout
	DryRun bool
	Error  error // The first error encountered.  If not nil, do not execute anything else.
}

type Cmd struct {
	Cmd         *exec.Cmd
	Name        string
	Args        []string
	inputCmd    *Cmd
	mergeStderr bool
	script      *Script
}

/** Run a command or pipeline */
func (t *Cmd) Run() {
	if t.script.DryRun {
		return
	}
	if t.script.HasError() {
		return
	}
	if t.Cmd.Stdout == nil {
		t.Cmd.Stdout = os.Stdout
	}
	if t.Cmd.Stderr == nil {
		if t.mergeStderr {
			t.Cmd.Stderr = t.Cmd.Stdout
		} else {
			t.Cmd.Stderr = os.Stderr
		}
	}
	var commands []*Cmd
	for c := t; c != nil; c = c.inputCmd {
		if t.Cmd.Stderr == nil {
			t.Cmd.Stderr = os.Stderr
		}
		if c.inputCmd != nil {
			var err error
			c.Cmd.Stdin, err = c.inputCmd.Cmd.StdoutPipe()
			if t.script.InError(err) {
				return
			}
		} else if c.Cmd.Stdin == nil {
			c.Cmd.Stdin = os.Stdin
		}
		commands = append(commands, c)
	}
	/* reverse the order */
	for i, j := 0, len(commands)-1; i < j; i, j = i+1, j-1 {
		commands[i], commands[j] = commands[j], commands[i]
	}

	for _, c := range commands {
		if t.script.InError(c.Cmd.Start()) {
			if t.script.Trace {
				fmt.Fprintln(os.Stderr, "start error", c.Name, c.Args)
			}
			return
		}
	}
	for _, c := range commands {
		if t.script.InError(c.Cmd.Wait()) {
			if t.script.Trace {
				fmt.Fprintln(os.Stderr, "wait error", c.Name, c.Args)
			}
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

/** Run and ignore the output.  Return success or failure. */
func (t *Cmd) ToNull() bool {
	t.ToWriter(&NullWriter{})
	return t.Script().Error == nil
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
func (t *Script) AddError(err error) {
	if t.Error == nil {
		t.Error = err
	}
}

/** Check if the argument is an error, or whether the script already has an error.
  If the argument is an error and it is the first error encountered, it becomes the script's error.
  Return true if the script has an error, either the given error or a previous one.  */
func (t *Script) InError(err error) bool {
	t.AddError(err)
	return t.HasError()
}

func (t *Cmd) Print() {
	args := make([]interface{}, 1+len(t.Args))
	args[0] = t.Name
	for i, arg := range t.Args {
		args[1+i] = arg
	}
	fmt.Println(args...)
}

/** Create a command without running it.  The command can be executed or piped to another command. */
func (t *Script) Cmd(name string, args ...string) *Cmd {
	r := &Cmd{Name: name, Args: args}
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
		r.Print()
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

func (t *Cmd) MergeStderr() *Cmd {
	t.mergeStderr = true
	return t
}

func (t *Cmd) Script() *Script {
	return t.script
}
