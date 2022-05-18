package internal

import (
	"fmt"
	"log"
	"os"
)

func ConfigureLogger(path string) {
	if path == "" {
		log.SetOutput(os.Stdout)
		return
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	log.SetOutput(f)
}
