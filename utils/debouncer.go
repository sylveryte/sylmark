package utils

import "time"

type SylDebouncer struct {
	timeDelay time.Duration
	timer     *time.Timer
}

func NewSylDebouncer(dur time.Duration) *SylDebouncer {
	return &SylDebouncer{
		timeDelay: dur,
	}
}

func (s *SylDebouncer) Debounce(fn func()) {

	if s.timer != nil {
		s.timer.Stop()
	}

	s.timer = time.AfterFunc(s.timeDelay, fn)

}
