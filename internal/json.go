package internal

import "encoding/json"

func FmtAsJson(comics []Comic) (string, error) {
	jsonBytes, err := json.Marshal(comics)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
