package main

import (
	"net/http"
	"prod23/internal/config"
	"prod23/internal/httpapi"
	"prod23/internal/logging"
	"prod23/internal/service"
	"prod23/internal/store"
)

func main() {
	c := config.Load()
	st, e := store.Open(c.DataDir)
	if e != nil {
		panic(e)
	}
	logging.New().Event("info", "starting", map[string]any{"addr": c.HTTPAddr})
	http.ListenAndServe(c.HTTPAddr, (&httpapi.Server{Svc: service.New(st)}).Handler())
}
