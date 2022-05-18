package internal

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestParsePage(t *testing.T) {
	t.Parallel()

	content, err := os.ReadFile("../tests/assets/page.html")
	html := string(content)
	if err != nil {
		t.Error(err)
	}

	t.Run("should find all comics", func(t *testing.T) {
		t.Parallel()

		comics, _ := parsePage(html)
		assert.Len(t, comics, 50)
	})

	t.Run("should extract the right max page information", func(t *testing.T) {
		t.Parallel()

		_, maxPageIndex := parsePage(html)
		assert.Equal(t, maxPageIndex, 9)
	})

	t.Run("should extract the right information", func(t *testing.T) {
		t.Parallel()

		comics, _ := parsePage(html)
		tests := []struct {
			index  int
			id     int
			title  string
			imgSrc string
			date   string
		}{
			{
				index:  0,
				id:     418,
				title:  "Am Anfang war das Wort",
				imgSrc: "https://img.zeit.de/administratives/kaenguru-comics/2022-05/14/original__ffffff",
				date:   "2022-05-14T05:00:11+02:00",
			},
			{
				index:  49,
				id:     369,
				title:  "Irgendwie Pech",
				imgSrc: "https://img.zeit.de/administratives/kaenguru-comics/2022-03/18/original__ffffff",
				date:   "2022-03-18T05:00:05+01:00",
			},
		}

		for _, test := range tests {
			assert.Equal(t, comics[test.index].Id, test.id)
			assert.Equal(t, comics[test.index].Title, test.title)
			assert.Equal(t, comics[test.index].ImgSrc, test.imgSrc)
			assert.Equal(t, comics[test.index].Date, test.date)
		}
	})
}
