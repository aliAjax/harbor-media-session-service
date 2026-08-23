package domain

func (s *Store) Release(lease Lease) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Only the owner that currently holds the lease (matching owner AND token)
	// may release it. Checking owner in addition to token prevents a node from
	// deleting another node's lease when it releases a stale handle that happens
	// to carry a matching token.
	if old, ok := s.leases[lease.RoomID]; ok && old.Token == lease.Token && old.Owner == lease.Owner {
		delete(s.leases, lease.RoomID)
	}
}
