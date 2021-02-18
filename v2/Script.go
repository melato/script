package script

import (
	"fmt"
	"os/exec"
)

/**
  Script is used to run and/or print external programs, using exec.Cmd
- Redirect output to file
- Copy output to a string or []string

- It maintains an internal error, and stops running or printing commands if it has an error.
*/

type Script struct {
	Trace  bool // print exec arguments to stdout
	DryRun bool
	Errors Errors
}

/** Return true if an error has happened */
func (t *Script) HasError() bool {
	return t.Errors.HasError()
}

func (t *Script) Error() error {
	return t.Errors.Error()
}

/** Check if the argument is an error, or whether the script already has an error.
  If the argument is an error and it is the first error encountered, it becomes the script's error.
  Return true if the script has an error, either the given error or a previous one.  */
func (t *Script) AddError(err error) {
	t.Errors.Handle(err)
}

/** Cmd creates a Cmd wrapper. */
func (t *Script) Command(name string, args ...string) *Cmd {
	return &Cmd{Cmd: exec.Command(name, args...), Script: t}
}

/** Create a command and run it or print it */
func (t *Script) Run(name string, args ...string) {
	if t.HasError() {
		return
	}
	t.RunCmd(exec.Command(name, args...))
}

/** Run or print zero or more commands */
func (t *Script) RunCmd(cmd ...*exec.Cmd) {
	if t.Trace {
		for i, c := range cmd {
			var prefix string
			if i > 0 {
				prefix = "| "
			}
			fmt.Printf("%s%s\n", prefix, c.String())
		}
	}
	if !t.DryRun {
		t.AddError(Run(cmd...))
	}
}
