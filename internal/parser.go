package internal

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

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

		width, height, ratio, err := FetchImageDimentions(imgSrc)
		if err != nil {
			log.Printf("[ERROR] Could not determine width/height for image URL %s", imgSrc)
		}

		comic := Comic{id, title, date, ComicImg{imgSrc, width, height, ratio}}
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
