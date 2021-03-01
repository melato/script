package script

import (
	"os/exec"
)

/**
  Script is used to run script.Cmd which are convenience wrappers around exec.Cmd
  Deprecated:  Use the lighterweight Run(), which uses exec.Cmd directly.
  scripting utilities:
- Run programs
- Run pipelines
- Redirect output to file
- Copy output to a string or []string

- Maintains an internal error state, so the user does not have to check for an error after running each command.
*/

type Script struct {
	Trace  bool // print exec arguments to stdout
	DryRun bool
	Error  error // The first error encountered.  If not nil, do not execute anything else.
}

/** Check if the argument is an error, or whether the script already has an error.
  If the argument is an error and it is the first error encountered, it becomes the script's error.
  Return true if the script has an error, either the given error or a previous one.  */
func (t *Script) InError(err error) bool {
	t.AddError(err)
	return t.HasError()
}

func (t *Script) Command(name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)
	if t.Trace {
		println("", cmd)
	}
	return cmd
}
