package internal

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"time"
)

const (
	urlTemplate = "https://www.zeit.de/serie/die-kaenguru-comics?p="
	startUrl    = urlTemplate + "1"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	httpClient HttpClient
)

func init() {
	httpClient = &http.Client{}
}

func fetchPageBody(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error when creating request with URL: %s", url)
		return "", err
	}

	// Zeit.de only allows direct access for certain user agents.
	// Set the user agent manually to avoid running into a ads/permission wall.
	req.Header.Set("User-Agent", "curl/7.82.0")

	res, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Error when fetching URL: %s", url)
		return "", err
	}

	if res.StatusCode >= 400 {
		return "", fmt.Errorf("issue while fetching the page, status code: %v", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error when reading server response")
		return "", err

	}
	res.Body.Close()

	return string(body), nil
}

// fetchAndExtract fetches a URL and extracts all comics from the body.
// It returns a slice of comics and the maximum page value that has been read
// from the paginator.
//
// The function can be used either by providing a channel and later reading from
// it, or by evaluating the returned values.
func fetchAndExtract(url string) ([]Comic, int) {
	html, err := fetchPageBody(url)
	if err != nil {
		log.Fatalf("Error when fetching URL: %s", url)
		panic(err)
	}
	comics, lastPageIndex := parsePage(html)
	log.Printf("Fetched %s\n", url)
	return comics, lastPageIndex
}

// FetchAll fetches the initial page and all found subpages to create a slice of
// all published comics. Comics are fetched from all subpages in parallel and
// are sorted after their ID before returned.
func FetchAll(timeout int) []Comic {
	comics, maxPageIndex := fetchAndExtract(startUrl)

	ch := make(chan []Comic)

	// looping from i=2 because we have already indexed the first page
	for i := 2; i <= maxPageIndex; i++ {
		url := urlTemplate + fmt.Sprintf("%v", i)
		go func() {
			cs, _ := fetchAndExtract(url)
			ch <- cs
		}()
	}

	// looping from i=2 because we have already indexed the first page
	for i := 2; i <= maxPageIndex; i++ {
		select {
		case next := <-ch:
			comics = append(comics, next[:]...)
		case <-time.After(time.Second * time.Duration(timeout)):
			continue
		}
	}

	// Comics have been received async and are out of order. Sort them now.
	sort.Slice(comics, func(a, b int) bool { return comics[a].Id < comics[b].Id })

	return comics
}
