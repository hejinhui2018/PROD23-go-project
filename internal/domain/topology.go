package domain

type Topology struct {
	FeederSections map[string][]string
	DeviceSections map[string]string
}

func (t Topology) Clone() Topology {
	n := Topology{FeederSections: map[string][]string{}, DeviceSections: map[string]string{}}
	for k, v := range t.FeederSections {
		n.FeederSections[k] = append([]string(nil), v...)
	}
	for k, v := range t.DeviceSections {
		n.DeviceSections[k] = v
	}
	return n
}
func (t Topology) FeederExists(f string) bool { _, ok := t.FeederSections[f]; return ok }
func (t Topology) DeviceExists(d string) bool { _, ok := t.DeviceSections[d]; return ok }
func (t Topology) SectionCount(f string) int  { return len(t.FeederSections[f]) }
func (t Topology) DevicesOn(f string) []string {
	out := []string{}
	for d, s := range t.DeviceSections {
		for _, x := range t.FeederSections[f] {
			if s == x {
				out = append(out, d)
			}
		}
	}
	return out
}

func DefaultTopology() Topology {
	return Topology{map[string][]string{"FDR-1": {"S1", "S2", "S3"}}, map[string]string{"SW-1": "S2", "SW-2": "S3"}}
}
func (t Topology) Affected(feeder, device string) []string {
	if s := t.DeviceSections[device]; s != "" {
		return []string{s}
	}
	return append([]string(nil), t.FeederSections[feeder]...)
}
