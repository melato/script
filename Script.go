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
	"bufio"
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
	IsRun       bool
	mergeStderr bool
	script      *Script
}

func (t *Cmd) Run() {
	if t.IsRun {
		return
	}
	if t.mergeStderr {
		t.Cmd.Stderr = t.Cmd.Stdout
	}
	t.script.checkError(t.Cmd.Run())
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

func (t *Cmd) Pipe(x *Cmd) {

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
	r.Cmd.Stdin = os.Stdin
	r.Cmd.Stdout = os.Stdout
	r.Cmd.Stderr = os.Stderr
	return r
}

func (t *Cmd) MergeStderr() {
	t.mergeStderr = true
}

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
