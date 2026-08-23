package domain

func validateTransitionTarget(j *Job, _ Status) error {
	if j == nil {
		return ErrInvalid
	}
	return nil
}
