package internal

import "os"

func ReadFileOrStdin(file string) (string, error) {
	if file == "-" {
		var content []byte
		_, err := os.Stdin.Read(content)
		if err != nil {
			return "", err
		}
		return string(content), nil
	} else {
		content, err := os.ReadFile(file)
		if err != nil {
			return "", err
		}
		return string(content), nil
	}
}
