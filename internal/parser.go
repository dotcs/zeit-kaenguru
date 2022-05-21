package internal

import (
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// parsePage evaluates the content of an HTML page that contains a list of
// comics.
// It returns a slice of comics and the maximum page value that has been read
// from the paginator.
//
// Each page contains up to 50 comics which are matched all at once with a
// simple regex without parsing the whole HTML document tree.
func parsePage(body string, timeout int) ([]Comic, int) {
	comicRe := regexp.MustCompile("(?s)<img.*?class=\"zon-teaser-standard__media-item\".*?alt=\"Folge ([0-9]+): (.*?)\".*?src=\"(.*?)\".*?<time.*?datetime=\"(.*?)\"")
	matches := comicRe.FindAllStringSubmatch(body, -1)
	comics := make([]Comic, 0)

	ch := make(chan Comic)

	for _, match := range matches {
		cMatch := match
		go func() {
			id, err := strconv.Atoi(cMatch[1])
			if err != nil {
				log.Println("[Comic parser] Could not extract ID of comic.")
			}
			title := cMatch[2]
			date := cMatch[4]

			// Rewrite image path to get the maximal resolution and a white
			// background.
			imgParts := strings.Split(cMatch[3], "/")
			imgParts = append(imgParts[0:len(imgParts)-1], "original__ffffff")
			imgSrc := strings.Join(imgParts, "/")

			width, height, ratio, err := FetchImageDimentions(imgSrc)
			if err != nil {
				log.Printf("[ERROR] Could not determine width/height for image URL %s", imgSrc)
			}

			comic := Comic{id, title, date, ComicImg{imgSrc, width, height, ratio}}
			ch <- comic
		}()
	}

	for i := 0; i < len(matches); i++ {
		select {
		case next := <-ch:
			comics = append(comics, next)
		case <-time.After(time.Second * time.Duration(timeout)):
			log.Printf("[ERROR] Parse page timed out for element %v\n", i)
			continue
		}
	}

	// Sort comics as their order is undefined through async action
	sort.Slice(comics, func(a, b int) bool { return comics[a].Id < comics[b].Id })

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
