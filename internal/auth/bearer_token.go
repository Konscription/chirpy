package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	// Directly get the "Authorization" header (http.Header is case-insensitive)
	authHeaders, ok := headers["Authorization"]
	if !ok || len(authHeaders) == 0 {
		return "", errors.New("missing Authorization header")
	}
	authHeader := authHeaders[0]
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("invalid Authorization header format")
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	token = strings.TrimSpace(token)
	if token == "" {
		return "", errors.New("empty bearer token")
	}
	return token, nil
}
