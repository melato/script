package script

import (
	"errors"
	"fmt"
	"io"
)

// ErrorHandler allows a function to report errors, without returning them directly
// It is used for convenience in order to minimize checking error return values from function calls
type ErrorHandler interface {
	// Add an error.  Do nothing if err is nil.
	Add(err error)
	// Return true if processing should continue.  Should return true if there were no errors.
	// May also return true if there were prior errors, but processing should continue in order to check for more errors.
	Continue() bool
}

// Errors implements ErrorHandler by collecting errors in a list.
type Errors struct {
	// The list of non-nil errors that have been passed to Add()
	List []error
	// If AlwaysContinue is true, Continue() always return true.
	// In that case, examine List, or HasErrors() to find out if there were errors
	AlwaysContinue bool
	// Writer is an optional io.Writer to write errors to.
	Writer io.Writer
}

func (t *Errors) Add(err error) {
	if err == nil {
		return
	}
	t.List = append(t.List, err)
	if t.Writer != nil {
		fmt.Fprintln(t.Writer, err)
	}
}

func (t *Errors) Continue() bool {
	return len(t.List) == 0 || t.AlwaysContinue
}

func (t *Errors) HasError() bool {
	return len(t.List) > 0
}

func (t *Errors) Clear() {
	t.List = nil
}

// Return the first error
func (t *Errors) First() error {
	if len(t.List) > 0 {
		return t.List[0]
	}
	return nil
}

// Return nil if there are no errors.  Otherwise, return the first error, or a generic error.
func (t *Errors) Error() error {
	if len(t.List) == 0 {
		return nil
	}
	if !t.AlwaysContinue {
		return t.First()
	}
	return errors.New("There were errors")
}
