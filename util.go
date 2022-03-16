package script

import (
	"bufio"
	"bytes"
)

func BytesToLines(data []byte) []string {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func IterateLines(data []byte, f func(string) error) error {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		err := f(scanner.Text())
		if err != nil {
			return err
		}
	}
	return nil
}
