package domain

import "time"

func (s *Store) Renew(lease Lease, ttl time.Duration) (Lease, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	old, ok := s.leases[lease.RoomID]
	if !ok || old.Owner != lease.Owner || old.Token != lease.Token {
		return Lease{}, ErrLeaseLost
	}
	lease.ExpiresAt = time.Now().Add(ttl)
	s.leases[lease.RoomID] = lease
	return lease, nil
}
