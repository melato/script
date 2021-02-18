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

// Cmd wraps exec.Cmd and adds convenience methods for redirecting its output in various ways.
type Cmd struct {
	Cmd         *exec.Cmd
	Script      *Script
	stdin       input
	mergeStderr bool
}

/** Set Cmd directory */
func (t *Cmd) Dir(dir string) *Cmd {
	t.Cmd.Dir = dir
	return t
}

/** Set Cmd directory */
func (t *Cmd) InputString(text string) *Cmd {
	t.stdin = &stringInput{Text: text}
	return t
}

/** Set Cmd directory */
func (t *Cmd) InputFile(file string) *Cmd {
	t.stdin = &fileInput{Path: file}
	return t
}

func (t *Cmd) Run() {
	if t.Script.HasError() {
		return
	}
	if t.Script.Trace {
		var suffix string
		var text []string
		if t.stdin != nil {
			text = t.stdin.TraceStrings()
			if len(text) > 0 {
				suffix = text[0]
				text = text[1:]
			}
		}
		fmt.Printf("%s%s\n", t.Cmd.String(), suffix)
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
		t.Script.AddError(Run(t.Cmd))
	}
}

/** Run and redirect output to a Writer */
func (t *Cmd) ToWriter(out io.Writer) {
	if t.Script.HasError() {
		return
	}
	t.Cmd.Stdout = out
	if t.mergeStderr {
		t.Cmd.Stderr = out
	}
	t.Run()
}

/** Run and redirect output to a file */
func (t *Cmd) ToFile(file string) {
	if t.Script.HasError() {
		return
	}
	fmt.Printf(" > %s\n", file)
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

/** Run and return the output */
func (t *Cmd) ToBytes() []byte {
	var buf bytes.Buffer
	t.ToWriter(&buf)
	return buf.Bytes()
}

/** Run and ignore the output.  Return success or failure. */
func (t *Cmd) ToNull() {
	t.ToWriter(ioutil.Discard)
}

/** Run and return the output as a string */
func (t *Cmd) ToString() string {
	return strings.TrimSpace(string(t.ToBytes()))
}

/** Run and return the output as a []string */
func (t *Cmd) ToLines() []string {
	return BytesToLines(t.ToBytes())
}

func (t *Cmd) MergeStderr() *Cmd {
	t.mergeStderr = true
	return t
}
