package main

import (
	"flag"
	"os"

	"github.com/dotcs/zeit-kaenguru/internal"
)

func santiyCheck(gist_id string) {
	_, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		panic("Environemnt variable 'GITHUB_TOKEN' must be set. Abort.")
	}
	if gist_id == "" {
		os.Stderr.Write([]byte("Gist ID is undefined. Abort."))
		os.Exit(1)
	}
}

func main() {
	logfile := flag.String("logfile", "", "defines the path to the logfile")
	gist_id := flag.String("gist-id", "", "ID of the gist that should be updated")
	file := flag.String("file", "", "Path to file that contains the json or '-' to read from stdin")
	flag.Parse()

	internal.ConfigureLogger(*logfile)

	santiyCheck(*gist_id)

	content, err := internal.ReadFileOrStdin(*file)
	if err != nil {
		panic(err)
	}

	files := make(map[string]internal.GistFile, 0)
	files["comics.json"] = internal.GistFile{Content: content, Filename: "comics.json"}

	gist := internal.Gist{
		Description: "Alle Kaenguru Comics von zeit.de (https://www.zeit.de/serie/die-kaenguru-comics)",
		Files:       files,
	}

	token := os.Getenv("GITHUB_TOKEN")
	err = internal.UpdateGist(*gist_id, token, gist)
	if err != nil {
		panic(err)
	}
}
