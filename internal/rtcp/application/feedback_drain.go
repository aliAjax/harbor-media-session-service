package application

import "harbor-sfu.local/harbor-sfu/internal/rtcp/domain"

func (q *FeedbackQueue) Drain() []domain.Feedback {
	q.mu.Lock()
	defer q.mu.Unlock()
	o := q.items
	q.items = q.items[:0]
	return o
}
