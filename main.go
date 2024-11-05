package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
  fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    cfg.fileserverHits.Add(1)
    next.ServeHTTP(w, r)
  })
}

func (cfg *apiConfig) hits(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hits: %d", cfg.fileserverHits.Load())
}

func (cfg *apiConfig) reset(w http.ResponseWriter, r * http.Request) {
  cfg.fileserverHits.Swap(0)
}

func healthz (w http.ResponseWriter,r *http.Request) {
  w.Header().Set("Content-Type", "text/plain; charset=utf-8")
  w.WriteHeader(200)
  w.Write([]byte("OK"))
}

func main() {
  mux := http.NewServeMux()
  apiCfg := apiConfig{}

  mux.Handle("/app", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
  mux.Handle("/app/assets/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
  mux.Handle("GET /api/healthz", http.HandlerFunc(healthz))
  mux.Handle("GET /api/metrics", http.HandlerFunc(apiCfg.hits))
  mux.Handle("POST /api/reset", http.HandlerFunc(apiCfg.reset))

  server := &http.Server {
    Addr: ":8080",
    Handler: mux,
  }
  if err := server.ListenAndServe(); err != nil {
  }
}

