package application

import (
	"harbor-sfu.local/harbor-sfu/internal/cluster/domain"
	"time"
)

type Coordinator struct {
	store *domain.Store
	node  string
	ttl   time.Duration
}

func NewCoordinator(store *domain.Store, node string, ttl time.Duration) *Coordinator {
	return &Coordinator{store: store, node: node, ttl: ttl}
}
