package script

import (
	"bufio"
	"bytes"
)

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
