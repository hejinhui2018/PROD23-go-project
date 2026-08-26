package cli

import (
	"bytes"
	"net/http"
)

func Post(url, payload string) (*http.Response, error) {
	return http.Post(url, "application/json", bytes.NewBufferString(payload))
}
