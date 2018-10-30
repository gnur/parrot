package main

import (
	"log/syslog"
	"net/http"
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
	for {
		log.WithFields(log.Fields{
			"source":  "generator",
			"date":    time.Now(),
			"counter": counter,
		}).Info("test")
		counter++
		log.WithFields(log.Fields{
			"source":  "generator",
			"date":    time.Now(),
			"counter": counter,
		}).Warning("test-warning")
		counter++
		time.Sleep(10 * time.Second)
	}

	log.Info("booksing is now running")
	log.Fatal(http.ListenAndServe(":7132", nil))
}
