package httpapi

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

func requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("x-request-id")
		if id == "" {
			b := make([]byte, 8)
			_, _ = rand.Read(b)
			id = hex.EncodeToString(b)
		}
		w.Header().Set("x-request-id", id)
		next.ServeHTTP(w, r)
	})
}
