package domain

import "time"

func (s *Store) Renew(lease Lease, ttl time.Duration) (Lease, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	old, ok := s.leases[lease.RoomID]
	// A lease may only be renewed while it is still valid: it must exist, still
	// belong to the same owner, carry the same token, and not have expired. Once
	// the TTL lapses the lease is lost and any holder must re-acquire rather than
	// resurrect a stale entry — otherwise an expired owner can silently revive a
	// lease another node may be racing to claim.
	if !ok || old.Owner != lease.Owner || old.Token != lease.Token || !old.ExpiresAt.After(now) {
		return Lease{}, ErrLeaseLost
	}
	lease.ExpiresAt = now.Add(ttl)
	s.leases[lease.RoomID] = lease
	return lease, nil
}
