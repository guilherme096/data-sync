package api

import (
	"log"
	"net/http"
	"time"

	"github.com/guilherme096/data-sync/internal/api/routers"
	"github.com/guilherme096/data-sync/pkg/data-sync/chatbot"
	datasync "github.com/guilherme096/data-sync/pkg/data-sync"
	"github.com/guilherme096/data-sync/pkg/data-sync/discovery"
	"github.com/guilherme096/data-sync/pkg/data-sync/matching"
	"github.com/guilherme096/data-sync/pkg/data-sync/query"
	"github.com/guilherme096/data-sync/pkg/data-sync/storage"
	"github.com/guilherme096/data-sync/pkg/data-sync/sync"
)

type Server struct {
	addr       string
	engine     datasync.QueryEngine
	storage    storage.MetadataStorage
	sync       sync.MetadataSync
	discovery  discovery.MetadataDiscovery
	agent      chatbot.AgentActions
	translator query.QueryTranslator
	matcher    *matching.Matcher
}

func NewServer(addr string, engine datasync.QueryEngine, storage storage.MetadataStorage, sync sync.MetadataSync, discovery discovery.MetadataDiscovery, agent chatbot.AgentActions, translator query.QueryTranslator, matcher *matching.Matcher) *Server {
	return &Server{
		addr:       addr,
		engine:     engine,
		storage:    storage,
		sync:       sync,
		discovery:  discovery,
		agent:      agent,
		translator: translator,
		matcher:    matcher,
	}
}

// corsMiddleware adds CORS headers to allow cross-origin requests
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
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

	globalRouter := routers.NewGlobalRouter(s.storage)
	globalRouter.RegisterRoutes(mux)

	relationRouter := routers.NewRelationRouter(s.storage, s.discovery, s.matcher)
	relationRouter.RegisterRoutes(mux)

	chatbotRouter := routers.NewChatbotRouter(s.agent, s.translator, s.discovery, s.storage)
	chatbotRouter.RegisterRoutes(mux)

	globalQueryRouter := routers.NewGlobalQueryRouter(s.translator)
	globalQueryRouter.RegisterRoutes(mux)

	// Wrap with CORS middleware
	handler := corsMiddleware(mux)

	server := &http.Server{
		Addr:         s.addr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Server starting on %s", s.addr)
	return server.ListenAndServe()
}
