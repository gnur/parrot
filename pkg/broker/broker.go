// Package broker starts up the syslog and web servers.
// broker also manages the relay of log
// messages from the syslog listener to the sse server.
package broker

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"time"

	"github.com/gnur/parrot/pkg/sse"
	"github.com/gnur/parrot/pkg/syslog"
	"github.com/gnur/parrot/pkg/webserver"
)

var reJSONExtract = regexp.MustCompile(`{.*}`)

// Config stores all options necessary for the entire workflow of
// where to receive syslog messages, whether to forward them, and
// where to serve the web dashboard.
type Config struct {
	Listeners URLs
	Silent    bool
	Web       string
}

// URLs satisfies the flag.Value interface for a list of URI
type URLs []*url.URL

func (us *URLs) String() string {
	return ""
}

// Set values or something
func (us *URLs) Set(value string) error {
	u, err := url.Parse(value)

	if err != nil {
		return err
	}

	*us = append(*us, u)

	return nil
}

// Start sanity-checks the configuration, creates and starts the servers
// with appropriate listeners, and beings forwarding syslog messages to
// the sse server.
func Start(c *Config) error {
	if !c.Silent {
		fmt.Println(` ______
(_____ \                     _
 _____) )___  ____ ____ ___ | |_
|  ____/ _  |/ ___) ___) _ \|  _)
| |   ( ( | | |  | |  | |_| | |__
|_|    \_||_|_|  |_|   \___/ \___))  v1.0.2`)
	}

	// At least one listener must be specified.
	if len(c.Listeners) == 0 {
		return errors.New("no mode specified! (TCP, UDP, or Unix required)")
	}

	// Create our myriad of servers.
	se := sse.New()
	sl := syslog.New()
	ws := webserver.New(c.Web, se)

	// Configure syslog server listeners.
	for _, u := range c.Listeners {
		switch u.Scheme {
		case "tcp":
			if err := sl.ListenTCP(u.Host); err != nil {
				return err
			}

			if !c.Silent {
				log.Printf("Listening for syslog messages at tcp://%s\n", u.Host)
			}
		case "udp":
			if err := sl.ListenUDP(u.Host); err != nil {
				return err
			}

			if !c.Silent {
				log.Printf("Listening for syslog messages at udp://%s\n", u.Host)
			}
		case "unix":
			if err := sl.ListenUnix(u.Path); err != nil {
				return err
			}

			if !c.Silent {
				log.Printf("Listening for syslog messages at unix://%s\n", u.Path)
			}
		default:
			return errors.New("invalid schema specified")
		}
	}

	// Start syslog receiving server
	if err := sl.Start(); err != nil {
		return err
	}

	// Take syslog messages and forward to SSE server
	// for dissemination to dashboard clients.
	go func() {
		for l := range sl.ReceiveLog {
			// These fields are of no interest to the dashboard.
			delete(l, "tls_peer")
			delete(l, "version")

			// Convert all timestamp to Unix epoch format, since they're
			// easier on networks and for Javascript to manipulate.
			if _, ok := l["timestamp"]; ok {
				l["timestamp"] = l["timestamp"].(time.Time).Unix()
			}
			var data []byte

			found := false
			val, ok := l["content"].(string)
			if ok {
				logContent := val
				data = []byte(reJSONExtract.FindString(logContent))
				found = true
			}
			// this means an logrus message was found
			if !found {

				var err error
				data, err = json.Marshal(l)
				if err != nil {
					// Silent failure
					return
				}
			}

			// Pass the transformed log to the SSE server.
			se.SendEvent <- &sse.Event{
				Event: "l",
				Data:  data,
			}
			ws.DataChan <- data
		}
	}()

	// And, finally, start the web server
	ws.Start()

	if !c.Silent {
		log.Printf("Parrot dashboard is available at http://%s\n", c.Web)
	}

	if !c.Silent {
		log.Println("Press Ctrl+C to stop.")
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Block here until Ctrl+C is caught.
	<-stop

	if !c.Silent {
		log.Println("Caught interrupt; shutting down cleanly...")
	}

	sl.Shutdown()
	ws.Shutdown()

	// Remove the Unix socket file(s) if we created any.
	for _, u := range c.Listeners {
		if u.Scheme == "unix" {
			if err := os.Remove(u.Path); err != nil {
				log.Printf("[ERROR] Failed to clean up unix socket at '%s'; you will need to remove it manually.", u.Path)
			}
		}
	}

	return nil
}
