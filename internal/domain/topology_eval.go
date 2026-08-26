package domain

import "sort"

type SectionRisk struct {
	Section  string
	Impact   int
	Isolated bool
}

func (t Topology) Evaluate(feeder, device string) []SectionRisk {
	sections := t.Affected(feeder, device)
	out := make([]SectionRisk, 0, len(sections))
	for i, s := range sections {
		out = append(out, SectionRisk{Section: s, Impact: len(sections) - i, Isolated: i > 0})
	}
	return out
}
func (t Topology) Validate() bool {
	for feeder, sections := range t.FeederSections {
		if feeder == "" || len(sections) == 0 {
			return false
		}
		seen := map[string]bool{}
		for _, s := range sections {
			if s == "" || seen[s] {
				return false
			}
			seen[s] = true
		}
	}
	return true
}
func (t Topology) Feeders() []string {
	r := make([]string, 0, len(t.FeederSections))
	for f := range t.FeederSections {
		r = append(r, f)
	}
	sort.Strings(r)
	return r
}
func (t Topology) Devices() []string {
	r := make([]string, 0, len(t.DeviceSections))
	for d := range t.DeviceSections {
		r = append(r, d)
	}
	sort.Strings(r)
	return r
}
func (t Topology) SectionExists(feeder, section string) bool {
	for _, s := range t.FeederSections[feeder] {
		if s == section {
			return true
		}
	}
	return false
}
