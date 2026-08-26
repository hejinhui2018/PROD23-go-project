package health

import (
	"os"
	"path/filepath"
)

func DataDirReady(dir string) bool {
	if dir == "" {
		return false
	}
	if e := os.MkdirAll(dir, 0755); e != nil {
		return false
	}
	f := filepath.Join(dir, ".healthcheck")
	if e := os.WriteFile(f, []byte("ok"), 0644); e != nil {
		return false
	}
	_ = os.Remove(f)
	return true
}
func Dependencies(events int, dir string) map[string]bool {
	return map[string]bool{"event_store": events >= 0, "data_dir": DataDirReady(dir)}
}
