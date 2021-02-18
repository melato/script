package script

import (
	"io"
	"os"
	"strings"
)

type input interface {
	Open() (io.Reader, error)
	Close() error
	// return zero or more strings that will be printed on trace.
	TraceStrings() []string
}

type stringInput struct {
	Text string
}

func (t *stringInput) Open() (io.Reader, error) {
	return strings.NewReader(t.Text), nil
}

func (t *stringInput) Close() error {
	return nil
}

func (t *stringInput) TraceStrings() []string {
	return []string{" << ---", t.Text, "---"}
}

type fileInput struct {
	Path string
	file *os.File
}

func (t *fileInput) Open() (io.Reader, error) {
	var err error
	t.file, err = os.Open(t.Path)
	return t.file, err
}

func (t *fileInput) Close() error {
	if t.file != nil {
		return t.file.Close()
	}
	return nil
}

func (t *fileInput) TraceStrings() []string {
	return []string{" < " + t.Path}
}
