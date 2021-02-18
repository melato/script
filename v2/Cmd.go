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
	mergeStderr bool
}

/** Set Cmd directory */
func (t *Cmd) Dir(dir string) *Cmd {
	t.Cmd.Dir = dir
	return t
}

func (t *Cmd) Run() {
	t.Script.RunCmd(t.Cmd)
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
