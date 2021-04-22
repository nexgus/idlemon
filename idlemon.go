package idlemon

import (
	"time"
)

type Monitor struct {
	Clear chan bool

	duration time.Duration
	timer    *time.Timer
	callback func(time.Time)
}

func NewMonitor(sec int64, f func(time.Time)) *Monitor {
	m := new(Monitor)
	m.Clear = make(chan bool, 1)
	m.callback = f
	m.duration = time.Duration(sec) * time.Second
	m.timer = time.NewTimer(m.duration)
	if ok := m.timer.Stop(); !ok {
		<-m.timer.C
	}

	return m
}

func (m *Monitor) Duration() time.Duration {
	return m.duration
}

func (m *Monitor) Run() {
	for {
		select {
		case <-m.Clear:
			if ok := m.timer.Stop(); !ok {
				// Ensure channel is drain and non-blocking
				select {
				case <-m.timer.C:
				default:
				}
			}
			m.timer.Reset(m.duration)
		case timeout := <-m.timer.C:
			if m.callback != nil {
				m.callback(timeout)
			}
		}
	}
}
