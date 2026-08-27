package domain

import "time"

type FaultStatus string

const (
	Assessed FaultStatus = "assessed"
	Restored FaultStatus = "restored"
	Review   FaultStatus = "manual_review"
)

type StepStatus string

const (
	Pending      StepStatus = "pending"
	Dispatched   StepStatus = "dispatched"
	Acknowledged StepStatus = "acknowledged"
	Completed    StepStatus = "completed"
	Failed       StepStatus = "failed"
)

type Fault struct {
	ID, Feeder, Device string
	Status             FaultStatus
	Sections           []string
	Version            int
	CreatedAt          time.Time
}
type Plan struct {
	ID, FaultID string
	Status      string
	Version     int
	Steps       []string
}
type Step struct {
	ID, PlanID, Section string
	Sequence            int
	Status              StepStatus
	ClaimedBy           string
	Version             int
	Idempotency         string
	Result              string
}
type Event struct {
	ID, Type    string
	AggregateID string
	Version     int
	Data        any
	At          time.Time
	Idempotency string
}
type ReviewRecord struct {
	ID, StepID, Reason, Worker string
	CreatedAt                  time.Time
	Resolved                   bool
}
