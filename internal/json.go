package internal

import "encoding/json"

func FmtAsJson(comics []Comic) (string, error) {
	jsonBytes, err := json.Marshal(comics)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func FmtGistAsJson(gist Gist) (string, error) {
	jsonBytes, err := json.Marshal(gist)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
