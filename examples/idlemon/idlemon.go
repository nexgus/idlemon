package main

import (
	"time"

	"github.com/nexgus/idlemon"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	timeout := make(chan time.Time, 1)
	monitor := idlemon.NewMonitor(
		5,
		timeout,
	)

	go monitor.Run()

	logrus.Debugf("clear/start")
	monitor.Clear <- true

	time.Sleep(2 * time.Second)
	logrus.Debugf("clear/start")
	monitor.Clear <- true

	t := <-timeout
	logrus.Debugf("timeout (%s)", t.Format(time.RFC3339Nano))
}
