package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dotcs/zeit-kaenguru/internal"
)

func main() {
	timeout := flag.Int("timeout", 60, "seconds until http requests time out")
	logfile := flag.String("logfile", "", "defines the path to the logfile")
	outputFile := flag.String("output-file", "", "defines where the result should be written to")
	flag.Parse()

	internal.ConfigureLogger(*logfile)

	comics := internal.FetchAll(*timeout)

	res, err := internal.FmtAsJson(comics)
	if err != nil {
		panic(err)
	}
	if *outputFile != "" {
		err = os.WriteFile(*outputFile, []byte(res), 0644)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println(res)
	}
	log.Println("👍 Success. Data fetched from server.")
}
