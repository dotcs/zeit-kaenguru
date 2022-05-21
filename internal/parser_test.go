package internal

import (
	"bytes"
	"io/ioutil"
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
	httpClientBak := httpClient

	// Mock http requests that fetch img sizes.
	// Return the same size (width=5613, height=2000) for every image based on
	// the real data for the comic with ID 1.
	mc := MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, req.Header.Get("Accept"), "image/webp")     // expect webp meme type
			assert.Equal(t, req.Header.Get("Range"), "bytes=0-64")      // expect limitation of requested file size
			content, _ := os.ReadFile("../tests/assets/1-64bytes.webp") // img size: 5613x2000
			r := ioutil.NopCloser(bytes.NewReader(content))
			return &http.Response{StatusCode: 200, Body: r}, nil
		},
	}
	httpClient = mc
	expectedImgWidth := 5613
	expectedImgHeight := 2000
	timeout := 10
	t.Cleanup(func() { httpClient = httpClientBak })

	content, err := os.ReadFile("../tests/assets/page.html")
	html := string(content)
	if err != nil {
		t.Error(err)
	}

	t.Run("should find all comics", func(t *testing.T) {
		t.Parallel()

		comics, _ := parsePage(html, timeout)
		assert.Len(t, comics, 50)
	})

	t.Run("should extract the right max page information", func(t *testing.T) {
		t.Parallel()

		_, maxPageIndex := parsePage(html, timeout)
		assert.Equal(t, maxPageIndex, 9)
	})

	t.Run("should extract the right information", func(t *testing.T) {
		t.Parallel()

		comics, _ := parsePage(html, timeout)

		type testImg struct {
			src    string
			height int
			width  int
		}
		type test struct {
			index int
			id    int
			title string
			date  string
			img   testImg
		}
		tests := []test{
			{
				index: 49,
				id:    418,
				title: "Am Anfang war das Wort",
				date:  "2022-05-14T05:00:11+02:00",
				img: testImg{
					src:    "https://img.zeit.de/administratives/kaenguru-comics/2022-05/14/original__ffffff",
					width:  expectedImgWidth,
					height: expectedImgHeight,
				},
			},
			{
				index: 0,
				id:    369,
				title: "Irgendwie Pech",
				date:  "2022-03-18T05:00:05+01:00",
				img: testImg{
					src:    "https://img.zeit.de/administratives/kaenguru-comics/2022-03/18/original__ffffff",
					width:  expectedImgWidth,
					height: expectedImgHeight,
				},
			},
		}

		for _, test := range tests {
			assert.Equal(t, test.id, comics[test.index].Id)
			assert.Equal(t, test.title, comics[test.index].Title)
			assert.Equal(t, test.date, comics[test.index].Date)
			assert.Equal(t, test.img.height, comics[test.index].Img.Height)
			assert.Equal(t, test.img.width, comics[test.index].Img.Width)
			assert.Equal(t, float32(test.img.width)/float32(test.img.height), comics[test.index].Img.Ratio)
			assert.Equal(t, test.img.src, comics[test.index].Img.Src)
		}
	})
}
