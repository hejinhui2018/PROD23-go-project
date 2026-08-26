package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

func call(method, url, key string, v any) (map[string]any, error) {
	b, _ := json.Marshal(v)
	r, e := http.NewRequest(method, url, bytes.NewReader(b))
	if e != nil {
		return nil, e
	}
	r.Header.Set("content-type", "application/json")
	if key != "" {
		r.Header.Set("Idempotency-Key", key)
	}
	x, e := http.DefaultClient.Do(r)
	if e != nil {
		return nil, e
	}
	defer x.Body.Close()
	raw, _ := io.ReadAll(x.Body)
	if x.StatusCode >= 300 {
		return nil, fmt.Errorf("%s: %s", x.Status, string(raw))
	}
	var out map[string]any
	_ = json.Unmarshal(raw, &out)
	return out, nil
}
func list(method, url string) ([]map[string]any, error) {
	r, e := http.NewRequest(method, url, nil)
	if e != nil {
		return nil, e
	}
	x, e := http.DefaultClient.Do(r)
	if e != nil {
		return nil, e
	}
	defer x.Body.Close()
	raw, _ := io.ReadAll(x.Body)
	if x.StatusCode >= 300 {
		return nil, fmt.Errorf("%s: %s", x.Status, string(raw))
	}
	var out []map[string]any
	if e := json.Unmarshal(raw, &out); e != nil {
		return nil, e
	}
	return out, nil
}
func main() {
	base := flag.String("base", "http://localhost:8080", "service URL")
	flag.Parse()
	f, e := call("POST", *base+"/faults", "smoke-fault", map[string]string{"feeder": "FDR-1", "device": "SW-1"})
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
	fid := f["ID"].(string)
	p, e := call("POST", *base+"/plans/"+fid, "smoke-plan", map[string]any{})
	if e != nil {
		panic(e)
	}
	steps := p["Steps"].([]any)
	for _, raw := range steps {
		id := raw.(string)
		for i, a := range []string{"claim", "ack", "complete"} {
			_, e = call("POST", *base+"/steps/"+id+"/"+a, "smoke-"+id+a, map[string]any{"worker": "smoke-crew", "version": i + 1, "result": "ok"})
			if e != nil {
				panic(e)
			}
		}
	}
	n, e := list("GET", *base+"/notifications")
	if e != nil {
		panic(e)
	}
	types := make([]string, 0, len(n))
	for _, item := range n {
		if value, ok := item["Type"].(string); ok {
			types = append(types, value)
		}
	}
	fmt.Printf("fault=%s plan=%s steps=%d notifications=%d types=%v\n", fid, p["ID"], len(steps), len(n), types)
}
