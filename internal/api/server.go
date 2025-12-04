package api

import (
	"log"
	"net/http"
	"time"

	datasync "github.com/guilherme096/data-sync/pkg/data-sync"
)

type Server struct {
	addr   string
	engine datasync.QueryEngine
}

func NewServer(addr string, engine datasync.QueryEngine) *Server {
	return &Server{
		addr:   addr,
		engine: engine,
	}
}

func (s *Server) Run() error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("POST /query", s.handleQuery)

	server := &http.Server{
		Addr:         s.addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Server starting on %s", s.addr)
	return server.ListenAndServe()
}
