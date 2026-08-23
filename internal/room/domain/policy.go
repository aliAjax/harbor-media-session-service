package domain

import (
	"fmt"
	"strings"
)

type Policy struct {
	AllowPublish   bool
	AllowSubscribe bool
	MaxTracks      int
	Codecs         []string
}

func (p Policy) ValidateTrack(t Track) error {
	if !p.AllowPublish {
		return fmt.Errorf("publishing disabled")
	}
	if p.MaxTracks == 0 {
		return fmt.Errorf("track limit is zero")
	}
	ok := len(p.Codecs) == 0
	for _, c := range p.Codecs {
		if strings.EqualFold(c, t.Codec) {
			ok = true
		}
	}
	if !ok {
		return fmt.Errorf("codec %s is not allowed", t.Codec)
	}
	return nil
}
func (p Policy) CanSubscribe() error {
	if !p.AllowSubscribe {
		return fmt.Errorf("subscriptions disabled")
	}
	return nil
}
