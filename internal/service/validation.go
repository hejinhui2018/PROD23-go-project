package service

import (
	"prod23/internal/domain"
	"strings"
)

func ValidateFaultInput(feeder, device string) error {
	if strings.TrimSpace(feeder) == "" || strings.TrimSpace(device) == "" {
		return domain.ErrInvalid
	}
	if len(feeder) > 64 || len(device) > 64 {
		return domain.ErrInvalid
	}
	return nil
}
func ValidateIdempotency(k string) error {
	if len(k) > 128 {
		return domain.ErrInvalid
	}
	if strings.ContainsAny(k, "\\r\\n") {
		return domain.ErrInvalid
	}
	return nil
}
func NormalizeResult(s string) string {
	r := strings.TrimSpace(s)
	if len(r) > 1024 {
		r = r[:1024]
	}
	return r
}
