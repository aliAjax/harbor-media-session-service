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
	return &FeedbackQueue{max: max, items: make([]domain.Feedback, 0, max)}
}
func (q *FeedbackQueue) Push(f domain.Feedback) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) >= q.max {
		q.items = q.items[1:]
	}
	q.items = append(q.items, f)
}
func (q *FeedbackQueue) Drain() []domain.Feedback {
	q.mu.Lock()
	defer q.mu.Unlock()
	o := append([]domain.Feedback(nil), q.items...)
	q.items = q.items[:0]
	return o
}
func (q *FeedbackQueue) Len() int { q.mu.Lock(); defer q.mu.Unlock(); return len(q.items) }
