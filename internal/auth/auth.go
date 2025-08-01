package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetApiKey(headers *http.Header) (string, error) {
	apiKey := headers.Get("Authorization")
	if apiKey == "" {
		return "", errors.New("no api key provided")
	}

	splitHeader := strings.Split(apiKey, " ")
	if len(splitHeader) != 2 || splitHeader[0] != "ApiKey" {
		return "", errors.New("invalid api key provided")
	}

	return splitHeader[1], nil
}
