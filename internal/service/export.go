package service

import (
	"encoding/json"
	"os"
	"path/filepath"
	"prod23/internal/domain"
)

type Export struct {
	Faults  []domain.Fault        `json:"faults"`
	Plans   []domain.Plan         `json:"plans"`
	Steps   []domain.Step         `json:"steps"`
	Reviews []domain.ReviewRecord `json:"reviews"`
}

func (s *Service) Export() Export {
	s.mu.RLock()
	defer s.mu.RUnlock()
	x := Export{Reviews: append([]domain.ReviewRecord(nil), s.Reviews...)}
	for _, f := range s.Faults {
		x.Faults = append(x.Faults, *f)
	}
	for _, p := range s.Plans {
		x.Plans = append(x.Plans, *p)
	}
	for _, st := range s.Steps {
		x.Steps = append(x.Steps, *st)
	}
	return x
}
func (s *Service) ExportFile(dir string) (string, error) {
	if e := os.MkdirAll(dir, 0755); e != nil {
		return "", e
	}
	p := filepath.Join(dir, "recovery-export.json")
	b, e := json.MarshalIndent(s.Export(), "", "  ")
	if e != nil {
		return "", e
	}
	if e = os.WriteFile(p, b, 0644); e != nil {
		return "", e
	}
	return p, nil
}
func ImportExport(b []byte) (Export, error) { var x Export; e := json.Unmarshal(b, &x); return x, e }
func (x Export) Valid() bool                { return len(x.Faults) >= 0 && len(x.Plans) >= 0 && len(x.Steps) >= 0 }
