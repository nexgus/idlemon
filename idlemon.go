package idlemon

import (
	"time"
)

type Monitor struct {
	Clear chan bool

	duration       time.Duration
	timer          *time.Timer
	isIdle         bool
	idleCallback   func(time.Time)
	resumeCallback func(time.Time)
}

func NewMonitor(sec int64, idleCB, resumeCB func(time.Time)) *Monitor {
	m := new(Monitor)
	m.Clear = make(chan bool, 1)
	m.isIdle = true
	m.idleCallback = idleCB
	m.resumeCallback = resumeCB
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
			if m.isIdle {
				if m.resumeCallback != nil {
					m.resumeCallback(time.Now())
				}
			}
			m.isIdle = false

			// Ensure channel is drain and non-blocking
			if ok := m.timer.Stop(); !ok {
				select {
				case <-m.timer.C:
				default:
				}
			}

			// Reset timer
			m.timer.Reset(m.duration)

		case timeout := <-m.timer.C:
			m.isIdle = true
			if m.idleCallback != nil {
				m.idleCallback(timeout)
			}
		}
	}
}
