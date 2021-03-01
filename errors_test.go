package script

import (
	"errors"
	"fmt"
	"testing"
)

func F(eh ErrorHandler, msg string) {
	if eh.Continue() {
		fmt.Println(msg)
	}
	eh.Add(errors.New("test"))
}

func TestUsage(t *testing.T) {
	errs := &Errors{}
	F(errs, "a")
	F(errs, "b")
}
