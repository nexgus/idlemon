package main

import (
	"time"

	"github.com/nexgus/idlemon"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.000000",
	})

	monitor := idlemon.NewMonitor(
		5,
		func(t time.Time) {
			logrus.Warnf("timeout: %s", t.Format(time.RFC3339Nano))
		},
	)

	go monitor.Run()

	logrus.Infof("clear/start")
	monitor.Clear <- true

	time.Sleep(2 * time.Second)
	logrus.Infof("clear/start")
	monitor.Clear <- true

	time.Sleep(monitor.Duration() + 500*time.Millisecond)
	logrus.Infof("stop")
}
