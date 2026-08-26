package domain

import (
	"encoding/json"
	"time"
)

type Envelope struct {
	Kind    string          `json:"kind"`
	Version int             `json:"version"`
	At      time.Time       `json:"at"`
	Payload json.RawMessage `json:"payload"`
}

func Encode(v any) (json.RawMessage, error) { return json.Marshal(v) }
func Decode(b []byte, v any) error          { return json.Unmarshal(b, v) }
func EventEnvelope(e Event) (Envelope, error) {
	b, err := json.Marshal(e.Data)
	return Envelope{Kind: e.Type, Version: e.Version, At: e.At, Payload: b}, err
}
func IsKnownFaultStatus(s FaultStatus) bool { return s == Assessed || s == Restored || s == Review }
func IsKnownStepStatus(s StepStatus) bool {
	return s == Pending || s == Dispatched || s == Acknowledged || s == Completed || s == Failed
}
