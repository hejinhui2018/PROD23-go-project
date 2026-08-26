package domain

func ValidateFault(f Fault) error {
	if f.ID == "" || f.Feeder == "" {
		return ErrInvalid
	}
	return nil
}
