package domain

func (j *Job) Transition(s Status) error {
	if err := validateTransitionTarget(j, s); err != nil {
		return err
	}
	if !canAdvanceRecording(j.Status, s) {
		return ErrInvalid
	}
	j.Status = s
	stampRecordingTransition(j, s)
	clearTransitionError(j)
	return nil
}
