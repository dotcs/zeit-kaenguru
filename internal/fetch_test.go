package internal

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchPageBody(t *testing.T) {
	t.Parallel()

	t.Run("should have curl user-agent", func(t *testing.T) {
		t.Parallel()

		mc := MockClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, req.Header.Get("User-Agent"), "curl/7.82.0")
				r := ioutil.NopCloser(bytes.NewReader([]byte("")))
				return &http.Response{StatusCode: 200, Body: r}, nil
			},
		}
		httpClient = mc
		fetchPageBody("https://example.com")
	})

	t.Run("should have curl user-agent", func(t *testing.T) {
		t.Parallel()

		mc := MockClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				r := ioutil.NopCloser(bytes.NewReader([]byte("the body")))
				return &http.Response{StatusCode: 200, Body: r}, nil
			},
		}
		httpClient = mc
		body, _ := fetchPageBody("https://example.com")
		assert.Equal(t, body, "the body")
	})

	t.Run("should raise an error if server has issue", func(t *testing.T) {
		t.Parallel()

		mc := MockClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				r := ioutil.NopCloser(bytes.NewReader([]byte("")))
				return &http.Response{StatusCode: 500, Body: r}, nil
			},
		}
		httpClient = mc
		_, err := fetchPageBody("https://example.com")
		if assert.Error(t, err) {
			assert.ErrorContains(t, err, "Status code: 500")
		}
	})
}
