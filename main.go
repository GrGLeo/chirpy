package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/GrGLeo/chirpy/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
  fileserverHits atomic.Int32
  dbQueries *database.Queries 
  platform string
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
  fmt.Printf("platform: %q", cfg.platform)
  cfg.fileserverHits.Swap(0)
  if cfg.platform != "dev" {
    http.Error(w, "Unauthorized", 403)
    return
  }
  err := cfg.dbQueries.DeleteUsers(r.Context())
  if err != nil {
    http.Error(w, "Error while deleting users", 500)
  }

  w.WriteHeader(http.StatusOK)
  w.Write([]byte("Success"))
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

func main() {
  godotenv.Load()
  dbUrl := os.Getenv("DB_URL")
  platform := os.Getenv("PLATFORM")
  
  db, _ := sql.Open("postgres", dbUrl)
  dbQueries := database.New(db)

  mux := http.NewServeMux()
  apiCfg := apiConfig{
    dbQueries: dbQueries,
    platform: platform,
  }

  mux.Handle("/app", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
  mux.Handle("/app/assets/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
  mux.Handle("GET /api/healthz", http.HandlerFunc(healthz))
  mux.Handle("GET /admin/metrics", http.HandlerFunc(apiCfg.metrics))
  mux.Handle("POST /admin/reset", http.HandlerFunc(apiCfg.reset))
  mux.Handle("POST /api/chirps", http.HandlerFunc(apiCfg.WriteChirps))
  mux.Handle("POST /api/users", http.HandlerFunc(apiCfg.CreateUser))
  mux.Handle("POST /api/login", http.HandlerFunc(apiCfg.UserLogin))
  mux.Handle("GET /api/chirps", http.HandlerFunc(apiCfg.GetChirps))
  mux.Handle("GET /api/chirps/{id}", http.HandlerFunc(apiCfg.GetChirp))

  server := &http.Server {
    Addr: ":8080",
    Handler: mux,
  }
  if err := server.ListenAndServe(); err != nil {
  }
}

