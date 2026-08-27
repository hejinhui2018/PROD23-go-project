package service

import "prod23/internal/domain"

func (s *Service) Audit() []domain.Event { return s.Store.Events() }
