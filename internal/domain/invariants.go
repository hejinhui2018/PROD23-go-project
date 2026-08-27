package domain

import "fmt"

func ValidatePlan(p Plan) error {
	if p.ID == "" || p.FaultID == "" || p.Status != "confirmed" || len(p.Steps) == 0 {
		return ErrInvalid
	}
	return nil
}
func ValidateStep(s Step) error {
	if s.ID == "" || s.PlanID == "" || s.Sequence < 1 || s.Version < 1 {
		return ErrInvalid
	}
	switch s.Status {
	case Pending, Dispatched, Acknowledged, Completed, Failed:
	default:
		return ErrInvalid
	}
	return nil
}
func ValidateWorker(worker string) error {
	if len(worker) < 2 || len(worker) > 80 {
		return fmt.Errorf("worker name length")
	}
	for _, r := range worker {
		if r == ' ' || r == '\n' || r == '\t' {
			return fmt.Errorf("worker name contains whitespace")
		}
	}
	return nil
}
func IsTerminal(s StepStatus) bool { return s == Completed || s == Failed }
func CanClaim(s Step) bool         { return s.Status == Pending && s.ClaimedBy == "" }
func CanAcknowledge(s Step, worker string) bool {
	return s.Status == Dispatched && s.ClaimedBy == worker
}
func CanFinish(s Step, worker string) bool { return s.Status == Acknowledged && s.ClaimedBy == worker }
