package internal

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
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
	log.Printf("Fetch %s ...\n", url)
	html, err := fetchPageBody(url)
	if err != nil {
		log.Fatalf("Error when fetching URL: %s", url)
		panic(err)
	}
	comics, lastPageIndex := parsePage(html)
	log.Printf("Finished fetching %s\n", url)
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

func FetchImageDimentions(url string) (int, int, float32, error) {
	log.Printf("Fetch image dimensions for %s\n", url)

	// First 64 bytes are enough to determine the image dimensions.
	// Do not read more than that to keep the load on the server and client as
	// low as possible.
	maxBytes := 64

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "image/webp")
	req.Header.Set("Range", fmt.Sprintf("bytes=0-%v", maxBytes))
	res, err := httpClient.Do(req)
	if err != nil {
		log.Printf("[ERROR] Could not read URL %s\n", url)
		return 0, 0, 0.0, err
	}

	reader := bufio.NewReader(res.Body)
	buf := make([]byte, maxBytes)
	reader.Read(buf)

	dir := os.TempDir()
	filename := filepath.Join(dir, "demo.webp")
	defer os.Remove(filename)
	ioutil.WriteFile(filename, buf, 0644)

	// Determine image dimensions through external `file` command and use regex
	// to extract the dimensions
	cmd := exec.Command("file", filename)
	stdout, err := cmd.Output()
	if err != nil {
		log.Printf("[ERROR] Could not execute external 'file' cmd for file %s\n", filename)
		return 0, 0, 0.0, err
	}

	re := regexp.MustCompile("([0-9]+)x([0-9]+)")
	matches := re.FindAllStringSubmatch(string(stdout), -1)
	if len(matches) == 0 {
		return 0, 0, 0.0, nil
	}
	width, err := strconv.Atoi(matches[0][1])
	if err != nil {
		log.Printf("[ERROR] Could not convert width value '%v' to int\n", width)
		return 0, 0, 0.0, err
	}
	height, _ := strconv.Atoi(matches[0][2])
	if err != nil {
		log.Printf("[ERROR] Could not convert height value '%v' to int\n", height)
		return 0, 0, 0.0, err
	}

	ratio := float32(width) / float32(height)

	return width, height, ratio, nil
}
