package application

import (
	"harbor-sfu.local/harbor-sfu/internal/rtcp/domain"
	"sync"
)

type FeedbackQueue struct {
	mu    sync.Mutex
	items []domain.Feedback
	max   int
}

func NewFeedbackQueue(max int) *FeedbackQueue {
	max = normalizeFeedbackMax(max)
	return &FeedbackQueue{max: max, items: make([]domain.Feedback, 0, max)}
}
