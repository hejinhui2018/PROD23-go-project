package health

import (
	"fmt"
	"time"
)

type Report struct {
	Status    string            `json:"status"`
	Events    int               `json:"events"`
	DataDir   string            `json:"data_dir"`
	CheckedAt time.Time         `json:"checked_at"`
	Checks    map[string]string `json:"checks"`
}

func Check(events int, dataDir string) Report {
	status := "ok"
	checks := map[string]string{"event_log": "ok", "topology": "ok", "notifications": "ok"}
	if events < 0 {
		status = "degraded"
		checks["event_log"] = "invalid"
	}
	return Report{Status: status, Events: events, DataDir: dataDir, CheckedAt: time.Now().UTC(), Checks: checks}
}
func Explain(r Report) string {
	if r.Status == "ok" {
		return "service ready"
	}
	return fmt.Sprintf("service %s", r.Status)
}
