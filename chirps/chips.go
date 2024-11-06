package chirps

import (
	"encoding/json"
	log"
	net/http"
	"strings"
)

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
      respondWihJson(w, 200, params.Body)
      }
  }

func respondWithError (w http.ResponseWriter, code int, msg string) {
}

func respondWihJson (w http.ResponseWriter, code int, msg string) {
  type returnVals struct {
    CleanedBody string `json:"cleaned_body"`
  }
  
  respBody := returnVals{
    CleanedBody: msg,
  }
  data, err := json.Marshal(respBody)
  if err != nil {
    log.Printf("Error marshaling json: %s", err)
    w.WriteHeader(500)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(code)
  w.Write(data)
}
  
func sanitizedChirp (msg string) (string, error) {
  profaneWord := [3]string{"kerfuffle", "sharbert", "fornax"}
  return "", nil
}
