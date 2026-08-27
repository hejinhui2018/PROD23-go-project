package service

func (s *Service) PendingSteps(plan string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	n := 0
	for _, x := range s.Steps {
		if x.PlanID == plan && x.Status == "pending" {
			n++
		}
	}
	return n
}
