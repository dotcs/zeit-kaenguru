package main

import (
	"flag"
	"fmt"

	"github.com/dotcs/zeit-kaenguru/internal"
)

func main() {
	timeout := flag.Int("timeout", 10, "seconds until http requests time out")
	logfile := flag.String("logfile", "", "defines the path to the logfile")
	flag.Parse()

	internal.ConfigureLogger(*logfile)

	comics := internal.FetchAll(*timeout)

	res, err := internal.FmtAsJson(comics)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
