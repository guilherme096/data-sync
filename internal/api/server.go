package api

import (
	"log"
	"net/http"
	"time"

	datasync "github.com/guilherme096/data-sync/pkg/data-sync"
	"github.com/guilherme096/data-sync/pkg/data-sync/discovery"
	"github.com/guilherme096/data-sync/pkg/data-sync/storage"
	"github.com/guilherme096/data-sync/pkg/data-sync/sync"
	"github.com/guilherme096/data-sync/internal/api/routers"
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

	// Register all routers
	healthRouter := routers.NewHealthRouter()
	healthRouter.RegisterRoutes(mux)

	queryRouter := routers.NewQueryRouter(s.engine)
	queryRouter.RegisterRoutes(mux)

	catalogsRouter := routers.NewCatalogsRouter(s.storage)
	catalogsRouter.RegisterRoutes(mux)

	discoveryRouter := routers.NewDiscoveryRouter(s.discovery)
	discoveryRouter.RegisterRoutes(mux)

	syncRouter := routers.NewSyncRouter(s.sync)
	syncRouter.RegisterRoutes(mux)

	server := &http.Server{
		Addr:         s.addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Server starting on %s", s.addr)
	return server.ListenAndServe()
}
