package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Workflow struct {
	Base   string
	Client *http.Client
}

func NewWorkflow(base string) *Workflow { return &Workflow{Base: base, Client: http.DefaultClient} }
func (w *Workflow) Do(method, path, key string, in, out any) error {
	b, _ := json.Marshal(in)
	r, e := http.NewRequest(method, w.Base+path, bytes.NewReader(b))
	if e != nil {
		return e
	}
	r.Header.Set("content-type", "application/json")
	if key != "" {
		r.Header.Set("Idempotency-Key", key)
	}
	resp, e := w.Client.Do(r)
	if e != nil {
		return e
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return fmt.Errorf("%s: %s", resp.Status, string(raw))
	}
	if out != nil {
		return json.Unmarshal(raw, out)
	}
	return nil
}
func (w *Workflow) Create(feeder, device string, out any) error {
	return w.Do("POST", "/faults", "cli-fault", map[string]string{"feeder": feeder, "device": device}, out)
}
func (w *Workflow) Confirm(id string, out any) error {
	return w.Do("POST", "/plans/"+id, "cli-plan", map[string]any{}, out)
}
