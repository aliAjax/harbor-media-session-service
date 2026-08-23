package application

import "time"

func nextReadDeadline() time.Time                  { return time.Now().Add(5 * time.Second) }
func deadlineExpired(now, deadline time.Time) bool { return !now.Before(deadline) }
