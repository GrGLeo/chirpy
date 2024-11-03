package main

import (
	"net/http"
)

func main() {
  mux := http.NewServeMux()

  mux.Handle("/", http.FileServer(http.Dir(".")))
  mux.Handle("/assets", http.FileServer(http.Dir(".")))
  mux.Handle("/healthz", http.HandlerFunc(healthz))
  server := &http.Server {
    Addr: ":8080",
    Handler: mux,
  }
  if err := server.ListenAndServe(); err != nil {
  }
}

func healthz (w http.ResponseWriter,r *http.Request) {
  w.Header().Set("Content-Type", "text/plain; charset=utf-8")
  w.WriteHeader(200)
  w.Write([]byte("ok"))
}
