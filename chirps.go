package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/GrGLeo/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
  ID        uuid.UUID       `json:"id"`
  CreatedAt time.Time       `json:"created_at"`
  UpdatedAt time.Time       `json:"updated_at"`
  Body      string          `json:"body"`
  UserId    uuid.NullUUID   `json:"user_id"`
}
  
type requestsBody struct {
  Body string `json:"body"`
  UserId uuid.NullUUID `json:"user_id"`
}

func (cfg *apiConfig) WriteChirps(w http.ResponseWriter, r *http.Request) {
  decoder := json.NewDecoder(r.Body)
  reqBody := requestsBody{}
  err := decoder.Decode(&reqBody)
  if err != nil {
    log.Printf("Error decoding parameters: %s", err)
    w.WriteHeader(500)
    return
  }

  if len(reqBody.Body) > 140 {
    http.Error(w, "Chirp is too long", 400)
    return
  }
  msg, _ := sanitizedChirp(reqBody.Body)
  fmt.Printf("this is the message: %q", msg)

  ChirpParam := database.CreateChirpParams{
    Body: msg,
    UserID: reqBody.UserId,
  }
  chirp, err := cfg.dbQueries.CreateChirp(r.Context(), ChirpParam)
  if err != nil {
    http.Error(w, "Internal server error", http.StatusInternalServerError)
    return
  }
  
  returnChirp := Chirp{
    ID: chirp.ID,
    CreatedAt: chirp.CreatedAt,
    UpdatedAt: chirp.UpdatedAt,
    Body: chirp.Body,
    UserId: chirp.UserID,
  }
   
  respondWihJson(w, 201, returnChirp)

  

}  

func respondWihJson (w http.ResponseWriter, code int, chirp Chirp) {
  data, err := json.Marshal(chirp)
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
  sentence := strings.Split(msg, " ")
  newSentence := []string{}
  for i, word := range sentence {
    changed := false
    for _, profane := range profaneWord {
      word = strings.ToLower(word)
      if strings.Contains(word, profane) {
        changed = true
        word = strings.ReplaceAll(word, profane, "****")
        newSentence = append(newSentence, word)
      }
    }
    if !changed {
      newSentence = append(newSentence, sentence[i])
    }
  }
  Sentence := strings.Join(newSentence, " ")
  return Sentence, nil
}

