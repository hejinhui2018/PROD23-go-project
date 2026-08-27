package httpapi

import (
	"errors"
	"net/http"
	"prod23/internal/domain"
)

func status(err error) int {
	if errors.Is(err, domain.ErrConflict) {
		return 409
	}
	if errors.Is(err, domain.ErrNotFound) {
		return 404
	}
	return 400
}
func writeError(w http.ResponseWriter, e error) { http.Error(w, e.Error(), status(e)) }
