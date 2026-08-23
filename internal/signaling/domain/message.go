package domain

type Message struct {
	Type          string         `json:"type"`
	ID            string         `json:"id"`
	Seq           uint64         `json:"seq"`
	RoomID        string         `json:"room_id,omitempty"`
	ParticipantID string         `json:"participant_id,omitempty"`
	TrackID       string         `json:"track_id,omitempty"`
	Payload       map[string]any `json:"payload,omitempty"`
}
type State struct {
	LastSeq       uint64
	Connected     bool
	RoomID        string
	ParticipantID string
}

func (s *State) Accept(m Message) bool {
	if m.Seq > 0 && m.Seq <= s.LastSeq {
		return false
	}
	if m.Seq > 0 {
		s.LastSeq = m.Seq
	}
	return true
}
