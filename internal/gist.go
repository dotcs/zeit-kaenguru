package internal

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Gist struct {
	Description string              `json:"description"`
	Files       map[string]GistFile `json:"files"`
}

type GistFile struct {
	Content  string `json:"content"`
	Filename string `json:"filename"`
}

func getGistUrl(id string) string {
	return fmt.Sprintf("https://api.github.com/gists/%s", id)
}

func UpdateGist(gistId string, ghToken string, data Gist) error {
	url := getGistUrl(gistId)
	req, err := http.NewRequest("PATCH", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", ghToken))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	body, err := FmtAsJson(data)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(strings.NewReader(body))

	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	println(fmt.Sprintf("%v", res.StatusCode))
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	res.Body.Close()

	return nil
}
