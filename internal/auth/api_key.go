package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	// Directly get the "Authorization" header (http.Header is case-insensitive)
	authHeaders, ok := headers["Authorization"]
	if !ok || len(authHeaders) == 0 {
		return "", errors.New("missing Authorization header")
	}
	authHeader := authHeaders[0]
	if !strings.HasPrefix(authHeader, "ApiKey ") {
		return "", errors.New("invalid Authorization header format")
	}
	token := strings.TrimPrefix(authHeader, "ApiKey ")
	token = strings.TrimSpace(token)
	if token == "" {
		return "", errors.New("empty ApiKey token")
	}
	return token, nil
}
