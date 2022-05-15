package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	urlTemplate = "https://www.zeit.de/serie/die-kaenguru-comics?p="
	startUrl    = urlTemplate + "1"
)

type Comic struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	ImgSrc string `json:"imgSrc"`
	Date   string `json:"date"`
}

func configureLogger(path string) {
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

func fetchPageBody(url string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	// Zeit.de only allows direct access for certain user agents.
	// Set the user agent manually to avoid running into a ads/permission wall.
	req.Header.Set("User-Agent", "curl/7.82.0")

	res, err := client.Do(req)
	if err != nil {
		log.Printf("Error when fetching URL: %s", url)
		return "", err
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	return fmt.Sprintf("%s", body), nil
}

// parsePage evaluates the content of an HTML page that contains a list of
// comics.
// It returns a slice of comics and the maximum page value that has been read
// from the paginator.
//
// Each page contains up to 50 comics which are matched all at once with a
// simple regex without parsing the whole HTML document tree.
func parsePage(body string) ([]Comic, int) {
	comicRe := regexp.MustCompile("(?s)<img.*?class=\"zon-teaser-standard__media-item\".*?alt=\"Folge ([0-9]+): (.*?)\".*?src=\"(.*?)\".*?<time.*?datetime=\"(.*?)\"")
	matches := comicRe.FindAllStringSubmatch(body, -1)
	comics := make([]Comic, 0)
	for _, match := range matches {
		id, err := strconv.Atoi(match[1])
		if err != nil {
			log.Println("[Comic parser] Could not extract ID of comic.")
		}
		title := match[2]
		date := match[4]

		// Rewrite image path to get the maximal resolution and a white
		// background.
		imgParts := strings.Split(match[3], "/")
		imgParts = append(imgParts[0:len(imgParts)-1], "original__ffffff")
		imgSrc := strings.Join(imgParts, "/")

		comic := Comic{id, title, imgSrc, date}
		comics = append(comics, comic)
	}

	pagesRe := regexp.MustCompile("(?s)<li class=\"pager__page\"><a href=\".*?([0-9+])\"")
	matchesPageIndex := pagesRe.FindAllStringSubmatch(body, -1)
	pageIndices := make([]int, 0)
	for _, match := range matchesPageIndex {
		index, err := strconv.Atoi(match[1])
		if err != nil {
			continue
		}
		pageIndices = append(pageIndices, index)
	}

	return comics, pageIndices[len(pageIndices)-1]
}

// fetchAndExtract fetches a URL, extracts all comics from the body and writes
// the slice in a channel if given.
// It returns a slice of comics and the maximum page value that has been read
// from the paginator.
//
// The function can be used either by providing a channel and later reading from
// it, or by evaluating the returned values.
func fetchAndExtract(url string, ch chan []Comic) ([]Comic, int) {
	html, err := fetchPageBody(url)
	if err != nil {
		log.Fatalf("Error when fetching URL: %s", url)
		panic(err)
	}
	comics, lastPageIndex := parsePage(html)
	if ch != nil {
		ch <- comics
	}
	log.Printf("Fetched %s\n", url)
	return comics, lastPageIndex
}

func fmtAsJson(comics []Comic) (string, error) {
	jsonBytes, err := json.Marshal(comics)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func main() {
	timeout := flag.Int("timeout", 10, "seconds until http requests time out")
	logfile := flag.String("logfile", "", "defines the path to the logfile")
	flag.Parse()

	configureLogger(*logfile)

	comics, maxPageIndex := fetchAndExtract(startUrl, nil)

	ch := make(chan []Comic)

	// looping from i=2 because we have already indexed the first page
	for i := 2; i <= maxPageIndex; i++ {
		url := urlTemplate + fmt.Sprintf("%v", i)
		go fetchAndExtract(url, ch)
	}

	// looping from i=2 because we have already indexed the first page
	for i := 2; i <= maxPageIndex; i++ {
		select {
		case next := <-ch:
			comics = append(comics, next[:]...)
		case <-time.After(time.Second * time.Duration(*timeout)):
			break
		}
	}

	// Comics have been received async and are out of order. Sort them now.
	sort.Slice(comics, func(a, b int) bool { return comics[a].Id < comics[b].Id })

	res, err := fmtAsJson(comics)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
