package domain

func (s *Store) Release(lease Lease) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if old, ok := s.leases[lease.RoomID]; ok && old.Token == lease.Token {
		delete(s.leases, lease.RoomID)
	}
}
