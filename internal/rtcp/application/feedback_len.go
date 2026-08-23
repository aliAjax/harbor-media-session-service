package application

func (q *FeedbackQueue) Len() int { q.mu.Lock(); defer q.mu.Unlock(); return len(q.items) }
