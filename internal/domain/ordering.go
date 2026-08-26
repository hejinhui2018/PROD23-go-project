package domain

import "sort"

func SortSteps(steps []Step) []Step {
	r := append([]Step(nil), steps...)
	sort.SliceStable(r, func(i, j int) bool {
		if r[i].Sequence == r[j].Sequence {
			return r[i].ID < r[j].ID
		}
		return r[i].Sequence < r[j].Sequence
	})
	return r
}
func NextStep(steps []Step) (Step, bool) {
	for _, s := range SortSteps(steps) {
		if s.Status == Pending {
			return s, true
		}
	}
	return Step{}, false
}
func CompletedCount(steps []Step) int {
	n := 0
	for _, s := range steps {
		if s.Status == Completed {
			n++
		}
	}
	return n
}
func HasFailure(steps []Step) bool {
	for _, s := range steps {
		if s.Status == Failed {
			return true
		}
	}
	return false
}
func PlanReady(steps []Step) bool { return len(steps) > 0 && CompletedCount(steps) == len(steps) }
