package store

import (
	"bufio"
	"encoding/json"
	"os"
	"prod23/internal/domain"
)

type RecoveryReport struct {
	Read, Skipped int
	LastVersion   map[string]int
}

func Scan(path string) ([]domain.Event, RecoveryReport, error) {
	f, e := os.Open(path)
	if e != nil {
		return nil, RecoveryReport{}, e
	}
	defer f.Close()
	r := RecoveryReport{LastVersion: map[string]int{}}
	out := []domain.Event{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		r.Read++
		var x domain.Event
		if json.Unmarshal(sc.Bytes(), &x) != nil {
			r.Skipped++
			continue
		}
		out = append(out, x)
		if x.Version > r.LastVersion[x.AggregateID] {
			r.LastVersion[x.AggregateID] = x.Version
		}
	}
	return out, r, sc.Err()
}
func ValidateSequence(events []domain.Event) bool {
	last := map[string]int{}
	for _, e := range events {
		if e.Version <= last[e.AggregateID] {
			return false
		}
		last[e.AggregateID] = e.Version
	}
	return true
}
