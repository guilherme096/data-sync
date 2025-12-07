package routers

import (
	"encoding/json"
	"net/http"

	"github.com/guilherme096/data-sync/pkg/data-sync/sync"
)

type SyncRouter struct {
	sync sync.MetadataSync
}

func NewSyncRouter(sync sync.MetadataSync) *SyncRouter {
	return &SyncRouter{
		sync: sync,
	}
}

func (r *SyncRouter) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /sync", r.handleSync)
}

func (r *SyncRouter) handleSync(w http.ResponseWriter, req *http.Request) {
	if err := r.sync.SyncAll(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Metadata sync completed successfully",
	})
}
