package service

import (
	"context"
	"prod23/internal/domain"
)

type Command struct {
	Action, StepID, Worker, Key, Result string
	Version                             int
}
type CommandResult struct {
	Step   *domain.Step
	Replay bool
}

func (s *Service) Execute(ctx context.Context, c Command) (CommandResult, error) {
	select {
	case <-ctx.Done():
		return CommandResult{}, ctx.Err()
	default:
	}
	var st *domain.Step
	var e error
	switch c.Action {
	case "claim":
		st, e = s.ClaimStep(c.StepID, c.Worker, c.Key, c.Version)
	case "ack":
		st, e = s.AcknowledgeStep(c.StepID, c.Worker, c.Key, c.Version)
	case "complete":
		st, e = s.CompleteStep(c.StepID, c.Worker, c.Key, c.Version, c.Result)
	case "fail":
		st, e = s.FailStep(c.StepID, c.Worker, c.Key, c.Version, c.Result)
	default:
		return CommandResult{}, domain.ErrInvalid
	}
	return CommandResult{Step: st}, e
}
func CommandActions() []string { return []string{"claim", "ack", "complete", "fail"} }
func IsCommandAction(a string) bool {
	for _, x := range CommandActions() {
		if x == a {
			return true
		}
	}
	return false
}
