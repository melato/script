package script

import (
	"bufio"
	"bytes"
	"io"
)

type NullWriter struct {
	io.Writer
}

func (t *NullWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
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
