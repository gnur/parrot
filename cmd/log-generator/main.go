package main

import (
	"log/syslog"
	"math/rand"
	"os"
	"time"

	lSyslog "github.com/sirupsen/logrus/hooks/syslog"

	log "github.com/sirupsen/logrus"
)

func main() {
	hook, err := lSyslog.NewSyslogHook("udp", "parrot:514", syslog.LOG_INFO, "")
	if err == nil {
		log.SetFormatter(&log.JSONFormatter{})
		log.AddHook(hook)
	}
	counter := 0
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "a host with no name"
	}
	for {
		rnd := rand.Intn(600) + 1
		time.Sleep(time.Duration(rnd) * time.Millisecond)

		log.WithFields(log.Fields{
			"source":  hostname,
			"counter": counter,
		}).Info("test")
		rnd = rand.Intn(600) + 1
		time.Sleep(time.Duration(rnd) * time.Millisecond)
		counter++
		log.WithFields(log.Fields{
			"source":  hostname,
			"counter": rnd,
		}).Warning("test-warning")
		counter++

	}
}
