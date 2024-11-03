package main

import (
  "fmt"
	"net/http"
)

func main() {
  mux := http.NewServeMux()

  mux.Handle("/", http.HandlerFunc(homepage))
  server := &http.Server {
    Addr: ":8080",
    Handler: mux,
  }
  if err := server.ListenAndServe(); err != nil {
  }
}


func homepage(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "Welcome to the homepage!")
}
