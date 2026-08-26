package domain

func ValidTransition(from, to StepStatus) bool {
	return (from == Pending && to == Dispatched) || (from == Dispatched && to == Acknowledged) || (from == Acknowledged && (to == Completed || to == Failed))
}
