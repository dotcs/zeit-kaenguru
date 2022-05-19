package internal

import (
	"bufio"
	"os"
)

var osStdin = os.Stdin

func ReadFileOrStdin(file string) (string, error) {
	if file == "-" {
		var stdin []byte

		scanner := bufio.NewScanner(osStdin)
		for scanner.Scan() {
			stdin = append(stdin, scanner.Bytes()...)
		}

		if err := scanner.Err(); err != nil {
			return "", err
		}

		return string(stdin), nil
	} else {
		content, err := os.ReadFile(file)
		if err != nil {
			return "", err
		}
		return string(content), nil
	}
}
