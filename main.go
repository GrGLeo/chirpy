package main

import (
  "fmt"
  "net/http"
  "sync/atomic"
  "log"
  "encoding/json"
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

func (cfg *apiConfig) metrics (w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html")
  w.WriteHeader(200)
  content := fmt.Sprintf(`
  <html>
  <body>
  <h1>Welcome, Chirpy Admin</h1>
  <p>Chirpy has been visited %d times!</p>
  </body>
  </html>
  `, cfg.fileserverHits.Load())
  w.Write([]byte(content))
}

func healthz (w http.ResponseWriter,r *http.Request) {
  w.Header().Set("Content-Type", "text/plain; charset=utf-8")
  w.WriteHeader(200)
  w.Write([]byte("OK"))
}

func validate_chirp(w http.ResponseWriter, r *http.Request) {
  type parameter struct {
    Body string `json:"body"`
  }

  decoder := json.NewDecoder(r.Body)
  params := parameter{}
  err := decoder.Decode(&params)
  if err != nil {
    log.Printf("Error decoding parameters: %s", err)
    w.WriteHeader(500)
    return
    }

    type returnVals struct {
      Valid bool `json:"valid"`
    }

    if len(params.Body) > 140 {
      respBody := returnVals{
        Valid: false,
      }
      data, err := json.Marshal(respBody)
      if err != nil {
        log.Printf("Error marshaling json: %s", err)
        w.WriteHeader(500)
        return
      }

      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(400)
      w.Write(data)
    } else {
      respBody := returnVals{
        Valid: true,
      }
      data, err := json.Marshal(respBody)
      if err != nil {
        log.Printf("Error marshaling json: %s", err)
        w.WriteHeader(500)
        return
      }
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(200)
      w.Write(data)
    }
  }


func main() {
  mux := http.NewServeMux()
  apiCfg := apiConfig{}

  mux.Handle("/app", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
  mux.Handle("/app/assets/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
  mux.Handle("GET /api/healthz", http.HandlerFunc(healthz))
  mux.Handle("GET /admin/metrics", http.HandlerFunc(apiCfg.metrics))
  mux.Handle("POST /admin/reset", http.HandlerFunc(apiCfg.reset))
  mux.Handle("POST /api/validate_chirp", http.HandlerFunc(validate_chirp))

  server := &http.Server {
    Addr: ":8080",
    Handler: mux,
  }
  if err := server.ListenAndServe(); err != nil {
  }
}

