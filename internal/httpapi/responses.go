package httpapi

import (
	"encoding/json"
	"net/http"
	"time"
)

type ErrorResponse struct {
	Error     string    `json:"error"`
	Code      int       `json:"code"`
	RequestID string    `json:"request_id,omitempty"`
	At        time.Time `json:"at"`
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
func writeMethod(w http.ResponseWriter, allow string) {
	w.Header().Set("allow", allow)
	writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed", Code: 405, At: time.Now().UTC()})
}
func writeAccepted(w http.ResponseWriter, v any) { writeJSON(w, http.StatusAccepted, v) }
func writeCreated(w http.ResponseWriter, v any)  { writeJSON(w, http.StatusCreated, v) }
func writeNoContent(w http.ResponseWriter)       { w.WriteHeader(http.StatusNoContent) }
