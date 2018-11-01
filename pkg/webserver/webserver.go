// Package webserver serves the web dashboard from either embedded files (prod)
// or app directory (dev), as well as provide the SSE server endpoint.
package webserver

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	// "path/filepath"
	"time"

	"github.com/gnur/parrot/pkg/sse"

	"github.com/GeertJohan/go.rice"
)

// Server is a http server
type Server struct {
	server   *http.Server
	DataChan chan json.RawMessage
	logCache logCache
}

type logCache struct {
	sync.Mutex
	logs []json.RawMessage
}

// New beings serving the web dashboard at the specified address, and
// sets up the endpoint for the provided SSE server.
func New(addr string, sseServer *sse.Server) *Server {
	r := http.NewServeMux()

	// Serves the dashboard app from either the local app directory in dev,
	// or embedded files in prod.
	// dir, _ := filepath.Abs("./app/dist")
	box := rice.MustFindBox("../../app/dist")

	webServer := &Server{
		server: &http.Server{
			Handler: r,
			Addr:    addr,
		},
		DataChan: make(chan json.RawMessage),
		logCache: logCache{},
	}
	r.Handle("/squawk", sseServer)
	r.HandleFunc("/logs", webServer.GetLogs())
	r.Handle("/", http.FileServer(box.HTTPBox()))

	return webServer
}

// Start begins listening and serving the dashboard & SSE endpoint.
func (s *Server) Start() {
	go s.server.ListenAndServe()
	go s.logAppender()
}

// GetLogs returns all logs that are in memory
func (s *Server) GetLogs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logCache.Lock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s.logCache.logs)
		s.logCache.Unlock()
		return
	}
}

func (s *Server) logAppender() {
	for l := range s.DataChan {
		s.logCache.Lock()
		s.logCache.logs = append(s.logCache.logs, l)
		if len(s.logCache.logs) > 2000 {
			s.logCache.logs = s.logCache.logs[1:]
		}
		s.logCache.Unlock()
	}

}

// Shutdown cleanly stops the server listener.
func (s *Server) Shutdown() error {
	ctx, can := context.WithTimeout(context.Background(), 1*time.Second)
	defer can()

	return s.server.Shutdown(ctx)
}
