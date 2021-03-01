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

// Cmd wraps exec.Cmd and adds convenience methods for redirecting stdin, stdout in various ways.
type Cmd struct {
	Cmd             *exec.Cmd
	Script          *Script
	stdin           input
	outputIndicator string
	combineOutput   bool
}

// Dir sets Cmd.Dir
func (t *Cmd) Dir(dir string) *Cmd {
	t.Cmd.Dir = dir
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

// Run runs and/or prints the command, applying any specified redirections
func (t *Cmd) Run() {
	if t.Script.HasError() {
		return
	}
	if t.Script.Trace {
		var inputIndicator string
		var text []string
		if t.stdin != nil {
			text = t.stdin.TraceStrings()
			if len(text) > 0 {
				inputIndicator = text[0]
				text = text[1:]
			}
		}
		fmt.Printf("%s%s%s\n", t.Cmd.String(), inputIndicator, t.outputIndicator)
		for _, s := range text {
			fmt.Println(s)
		}
	}
	if t.Script.DryRun {
		return
	}
	if t.stdin != nil {
		var err error
		t.Cmd.Stdin, err = t.stdin.Open()
		if err != nil {
			t.Script.AddError(err)
			return
		}
		defer t.stdin.Close()
	}
	if t.combineOutput {
		if t.Cmd.Stdout == nil {
			t.Cmd.Stdout = os.Stdout
		}
		t.Cmd.Stderr = t.Cmd.Stdout
	}
	t.Script.AddError(Run(t.Cmd))
}

// ToWriter redirects the output to an io.Writer and runs the command
func (t *Cmd) ToWriter(out io.Writer) {
	if t.Script.HasError() {
		return
	}
	t.Cmd.Stdout = out
	if t.combineOutput {
		t.Cmd.Stderr = out
	}
	t.Run()
}

// ToFile redirects the output to a file, runs the command, and closes the file
func (t *Cmd) ToFile(file string) {
	if t.Script.HasError() {
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
