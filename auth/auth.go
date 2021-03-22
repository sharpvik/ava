package auth

import "net/http"

const apiKeyHeaderName = "Ava-API-Key"

func Auth(apiKey string, r *http.Request) bool {
	return r.Header.Get(apiKeyHeaderName) == apiKey
}
