package api

import (
	"log"
	"net/http"
	"time"

	datasync "github.com/guilherme096/data-sync/pkg/data-sync"
	"github.com/guilherme096/data-sync/pkg/data-sync/discovery"
	"github.com/guilherme096/data-sync/pkg/data-sync/storage"
	"github.com/guilherme096/data-sync/pkg/data-sync/sync"
)

type Server struct {
	addr      string
	engine    datasync.QueryEngine
	storage   storage.MetadataStorage
	sync      sync.MetadataSync
	discovery discovery.MetadataDiscovery
}

func NewServer(addr string, engine datasync.QueryEngine, storage storage.MetadataStorage, sync sync.MetadataSync, discovery discovery.MetadataDiscovery) *Server {
	return &Server{
		addr:      addr,
		engine:    engine,
		storage:   storage,
		sync:      sync,
		discovery: discovery,
	}
}

func (s *Server) Run() error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("POST /query", s.handleQuery)

	// metadata
	mux.HandleFunc("GET /catalogs", s.handleListCatalogs)
	mux.HandleFunc("GET /catalogs/{name}", s.handleGetCatalog)
	mux.HandleFunc("GET /catalogs/{name}/schemas", s.handleListSchemas)
	mux.HandleFunc("POST /sync", s.handleSync)

	// discovery
	mux.HandleFunc("GET /discover/catalogs/{catalog}/schemas/{schema}/tables", s.handleDiscoverTables)
	mux.HandleFunc("GET /discover/catalogs/{catalog}/schemas/{schema}/tables/{table}/columns", s.handleDiscoverColumns)

	server := &http.Server{
		Addr:         s.addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Server starting on %s", s.addr)
	return server.ListenAndServe()
}
