package api

import (
	"log"
	"net/http"
	"time"
)

type Server struct {
	addr string
	// In the future, we will inject the 'core.Service' here
	// service core.Service
}

func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
	}
}

func (s *Server) Run() error {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("GET /health", s.handleHealth)

	server := &http.Server{
		Addr:         s.addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Server starting on %s", s.addr)
	return server.ListenAndServe()
}
