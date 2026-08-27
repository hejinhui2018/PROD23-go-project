package httpapi

import (
	"net/http/httptest"
	"prod23/internal/service"
	"prod23/internal/store"
	"strings"
	"testing"
)

func TestCreateFaultHTTP(t *testing.T) {
	st, _ := store.Open(t.TempDir())
	h := (&Server{Svc: service.New(st)}).Handler()
	r := httptest.NewRequest("POST", "/faults", strings.NewReader(`{"feeder":"FDR-1","device":"SW-1"}`))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 201 {
		t.Fatalf("status %d", w.Code)
	}
}
