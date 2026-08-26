package httpapi

import (
	"encoding/json"
	"net/http"
	"prod23/internal/domain"
	"prod23/internal/service"
	"strings"
)

type Server struct{ Svc *service.Service }

func (s *Server) Handler() http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("/healthz", s.health)
	m.HandleFunc("/metrics", s.metrics)
	m.HandleFunc("/faults", s.faults)
	m.HandleFunc("/faults/", s.fault)
	m.HandleFunc("/plans/", s.plans)
	m.HandleFunc("/steps/", s.steps)
	m.HandleFunc("/notifications", s.notifications)
	m.HandleFunc("/audit", s.audit)
	return requestID(m)
}
func write(w http.ResponseWriter, v any) {
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
func decode(r *http.Request, v any) error {
	if r.Body == nil {
		return domainErr("empty body")
	}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	return d.Decode(v)
}

type domainErr string

func (e domainErr) Error() string { return string(e) }
func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	write(w, map[string]any{"status": "ok", "events": len(s.Svc.Audit())})
}
func (s *Server) metrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/plain")
	w.Write([]byte("recovery_events_total " + itoa(len(s.Svc.Audit())) + "\n"))
}
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	b := []byte{}
	for i > 0 {
		b = append([]byte{byte('0' + i%10)}, b...)
		i /= 10
	}
	return string(b)
}
func (s *Server) faults(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	var q struct {
		Feeder string `json:"feeder"`
		Device string `json:"device"`
	}
	if e := decode(r, &q); e != nil {
		writeError(w, e)
		return
	}
	v, e := s.Svc.CreateFaultWithKey(q.Feeder, q.Device, r.Header.Get("Idempotency-Key"))
	if e != nil {
		writeError(w, e)
		return
	}
	w.WriteHeader(201)
	write(w, v)
}
func (s *Server) fault(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/faults/")
	if r.Method != "GET" {
		w.WriteHeader(405)
		return
	}
	v, e := s.Svc.GetFault(id)
	if e != nil {
		writeError(w, e)
		return
	}
	write(w, v)
}
func (s *Server) plans(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/plans/")
	if r.Method == "POST" {
		v, e := s.Svc.ConfirmWithKey(id, r.Header.Get("Idempotency-Key"))
		if e != nil {
			writeError(w, e)
			return
		}
		w.WriteHeader(201)
		write(w, v)
		return
	}
	if r.Method == "GET" {
		p, st, e := s.Svc.GetPlan(id)
		if e != nil {
			writeError(w, e)
			return
		}
		write(w, map[string]any{"plan": p, "steps": st})
		return
	}
	w.WriteHeader(405)
}
func (s *Server) steps(w http.ResponseWriter, r *http.Request) {
	x := strings.Split(strings.TrimPrefix(r.URL.Path, "/steps/"), "/")
	if len(x) != 2 {
		http.Error(w, "bad path", 400)
		return
	}
	var q struct {
		Worker  string `json:"worker"`
		Version int    `json:"version"`
		Result  string `json:"result"`
	}
	if e := decode(r, &q); e != nil {
		writeError(w, e)
		return
	}
	k := r.Header.Get("Idempotency-Key")
	var v *domain.Step
	var e error
	switch x[1] {
	case "claim":
		v, e = s.Svc.ClaimStep(x[0], q.Worker, k, q.Version)
	case "ack":
		v, e = s.Svc.AcknowledgeStep(x[0], q.Worker, k, q.Version)
	case "complete":
		v, e = s.Svc.CompleteStep(x[0], q.Worker, k, q.Version, q.Result)
	case "fail":
		v, e = s.Svc.FailStep(x[0], q.Worker, k, q.Version, q.Result)
	default:
		http.Error(w, "unknown action", 400)
		return
	}
	if e != nil {
		writeError(w, e)
		return
	}
	write(w, v)
}
func (s *Server) notifications(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(405)
		return
	}
	write(w, s.Svc.Queue.Pending())
}
func (s *Server) audit(w http.ResponseWriter, r *http.Request) {
	write(w, map[string]any{"events": s.Svc.Audit(), "reviews": s.Svc.ReviewsList()})
}
