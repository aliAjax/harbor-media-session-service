package application

import "harbor-sfu.local/harbor-sfu/internal/rtcp/domain"

func (q *FeedbackQueue) Push(f domain.Feedback) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) >= q.max {
		q.items = q.items[1:]
	}
	q.items = append(q.items, f)
}
