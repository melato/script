package script

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type Cmd struct {
	Cmd         *exec.Cmd
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

/** Set Cmd directory */
func (t *Cmd) Dir(dir string) *Cmd {
	t.Cmd.Dir = dir
	return t
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
	if t.script.DryRun {
		fmt.Printf(" > %s\n", file)
		return
	}
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
	t.ToWriter(ioutil.Discard)
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
	t.script = to.script // use the same script, the one we return with.
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
func (t *Cmd) Print() {
	Print(t.Cmd)
}

/** Create a command without running it.  The command can be executed or piped to another command. */
func (t *Script) Cmd(name string, args ...string) *Cmd {
	cmd := &Cmd{Cmd: exec.Command(name, args...)}
	if t.Trace {
		Print(cmd.Cmd)
	}
	return cmd
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
