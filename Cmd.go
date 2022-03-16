package script

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// Cmd wraps 0 or more exec.Cmd and adds convenience methods for redirecting stdin, stdout in various ways.
type Cmd struct {
	Commands        []*exec.Cmd
	Script          *Script
	stdin           input
	outputIndicator string
	combineOutput   bool
}

// Pipeline creates an empty pipeline
func (t *Script) Pipeline() *Cmd {
	return &Cmd{Script: t}
}

// Add - add an os/exec.Cmd to the command pipeline
func (t *Cmd) Add(cmd *exec.Cmd) *Cmd {
	t.Commands = append(t.Commands, cmd)
	return t
}

// PipeTo - shortcut to Add(exec.Command(name, arg...))
func (t *Cmd) PipeTo(name string, arg ...string) *Cmd {
	return t.Add(exec.Command(name, arg...))
}

// Dir sets Cmd.Dir for all commands so far.
func (t *Cmd) Dir(dir string) *Cmd {
	for _, cmd := range t.Commands {
		cmd.Dir = dir
	}
	return t
}

// InputString sets Cmd.Stdin to the given string input
func (t *Cmd) InputString(text string) *Cmd {
	t.stdin = &stringInput{Text: text}
	return t
}

// InputBytes sets Cmd.Stdin to the given []byte input
func (t *Cmd) InputBytes(bytes []byte) *Cmd {
	t.stdin = &bytesInput{Bytes: bytes}
	return t
}

// InputString sets Cmd.Stdin to the given file
func (t *Cmd) InputFile(file string) *Cmd {
	t.stdin = &fileInput{Path: file}
	return t
}

func (t *Cmd) validateCommands() bool {
	if t.Script.HasError() {
		return false
	}
	if len(t.Commands) > 0 {
		return true
	}
	t.Script.AddError(errors.New("no commands"))
	return false
}

// Run runs and/or prints the command, applying any specified redirections
func (t *Cmd) Run() {
	if !t.validateCommands() {
		return
	}
	if t.Script.Trace {
		var inputIndicator string
		var inputText []string
		if t.stdin != nil {
			inputText = t.stdin.TraceStrings()
			if len(inputText) > 0 {
				inputIndicator = inputText[0]
				inputText = inputText[1:]
			}
		}
		n := len(t.Commands)
		for i, cmd := range t.Commands {
			var out string
			if i == n-1 {
				out = t.outputIndicator
			}
			if i == 0 {
				fmt.Printf("%s%s%s\n", cmd.String(), inputIndicator, out)
				for _, s := range inputText {
					fmt.Println(s)
				}
			} else {
				fmt.Printf("| %s%s\n", cmd.String(), out)
			}
		}
	}
	if t.Script.DryRun {
		return
	}
	if t.stdin != nil {
		var err error
		t.Commands[0].Stdin, err = t.stdin.Open()
		if err != nil {
			t.Script.AddError(err)
			return
		}
		defer t.stdin.Close()
	}
	if t.combineOutput {
		cmd := t.Commands[len(t.Commands)-1]
		if cmd.Stdout == nil {
			cmd.Stdout = os.Stdout
		}
		cmd.Stderr = cmd.Stdout
	}
	t.Script.AddError(Run(t.Commands...))
}

// ToWriter redirects the output to an io.Writer and runs the command
func (t *Cmd) ToWriter(out io.Writer) {
	if !t.validateCommands() {
		return
	}
	cmd := t.Commands[len(t.Commands)-1]
	cmd.Stdout = out
	if t.combineOutput {
		cmd.Stderr = out
	}
	t.Run()
}

// ToFile redirects the output to a file, runs the command, and closes the file
func (t *Cmd) ToFile(file string) {
	if !t.validateCommands() {
		return
	}
	t.outputIndicator = " > " + file
	if t.Script.DryRun {
		return
	}
	f, err := os.Create(file)
	if err != nil {
		t.Script.AddError(err)
		return
	}
	defer f.Close()
	t.ToWriter(f)
}

// ToBytes runs the command and returns its stdout as a []byte
func (t *Cmd) ToBytes() []byte {
	var buf bytes.Buffer
	t.ToWriter(&buf)
	return buf.Bytes()
}

// ToNull runs the command and discards its output
func (t *Cmd) ToNull() {
	t.ToWriter(ioutil.Discard)
}

// ToString runs the command and returns its output as a string
func (t *Cmd) ToString() string {
	return strings.TrimSpace(string(t.ToBytes()))
}

// ToLines runs the command, splits its output into lines, and returns the lines
func (t *Cmd) ToLines() []string {
	return BytesToLines(t.ToBytes())
}

// CombineOutput redirects stderr to stdout.  If the output is returned or redirected to a file, it includes both stdout and stderr
func (t *Cmd) CombineOutput() *Cmd {
	t.combineOutput = true
	return t
}
