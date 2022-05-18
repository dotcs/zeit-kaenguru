package internal

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func ConfigureLogger(path string) {
	if path == "" {
		dir := os.TempDir()
		path = filepath.Join(dir, "log.log")
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	log.SetOutput(f)
}
