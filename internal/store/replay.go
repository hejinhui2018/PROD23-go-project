package store

func (s *Store) ReplayCount() int { return len(s.events) }
