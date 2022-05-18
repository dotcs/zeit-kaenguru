package internal

import "encoding/json"

func FmtAsJson[K any](obj K) (string, error) {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
