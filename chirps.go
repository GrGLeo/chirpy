package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/GrGLeo/chirpy/internal/auth"
	"github.com/GrGLeo/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
  ID        uuid.UUID       `json:"id"`
  CreatedAt time.Time       `json:"created_at"`
  UpdatedAt time.Time       `json:"updated_at"`
  Body      string          `json:"body"`
  UserId    uuid.UUID       `json:"user_id"`
}
  
type requestsBody struct {
  Body string `json:"body"`
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
  token, err := auth.GetBearerToken(r.Header)
  if err != nil {
    http.Error(w, err.Error(), http.StatusUnauthorized)
  }
  userId, err := auth.ValidateJWT(token, cfg.jwtsecret)
  if err != nil {
    http.Error(w, err.Error(), http.StatusUnauthorized)
  }

  if len(reqBody.Body) > 140 {
    http.Error(w, "Chirp is too long", 400)
    return
  }
  msg, _ := sanitizedChirp(reqBody.Body)
  fmt.Printf("this is the message: %q", msg)

  ChirpParam := database.CreateChirpParams{
    Body: msg,
    UserID: userId,
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
    UserId: userId,
  }
   
  respondWihJson(w, 201, returnChirp)
}  

func (cfg *apiConfig) GetChirps(w http.ResponseWriter, r *http.Request) {
  var AllChirps []Chirp
  var err error
  var Chirps []database.Chirp


  authorID := uuid.Nil
  userId := r.URL.Query().Get("author_id")
  if userId != "" {
    authorID, err = uuid.Parse(userId)
    if err != nil {
      http.Error(w, "Internal Server Error", http.StatusInternalServerError)
      return
    }
  }

  
  orderBy := r.URL.Query().Get("sort")
  if orderBy == "asc" {
    Chirps, err = cfg.dbQueries.GetChirpsAsc(r.Context())
  } else if orderBy == "desc" {
    Chirps, err = cfg.dbQueries.GetChirpsDesc(r.Context())
  } else {
    Chirps, err = cfg.dbQueries.GetChirpsAsc(r.Context())
  }

  if err != nil {
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    return
  }

  for _, result := range Chirps {
    if authorID != uuid.Nil && result.UserID != authorID {
      continue
    }
    AllChirps = append(
      AllChirps,
      Chirp{
        ID: result.ID,
        CreatedAt: result.CreatedAt,
        UpdatedAt: result.UpdatedAt,
        Body: result.Body,
        UserId: result.UserID,
      },
    )
  }
  respondWithJson(w, 200, AllChirps)
}

func (cfg *apiConfig) GetChirp(w http.ResponseWriter, r *http.Request) {
  id, err := uuid.Parse(r.PathValue("id"))
  if err != nil {
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    return
  }

  chirp, err := cfg.dbQueries.GetChirp(r.Context(), id)
  if err != nil {
    if err == sql.ErrNoRows {
      http.Error(w, "Chirp not found", http.StatusNotFound)
      return
    }
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  returnChirp := Chirp{
        ID: chirp.ID,
        CreatedAt: chirp.CreatedAt,
        UpdatedAt: chirp.UpdatedAt,
        Body: chirp.Body,
        UserId: chirp.UserID,
      }
  respondWihJson(w, 200, returnChirp)
}


func (cfg *apiConfig) DeleteChirp(w http.ResponseWriter, r *http.Request) {
  token, err := auth.GetBearerToken(r.Header)
  if err != nil {
    http.Error(w, err.Error(), http.StatusUnauthorized)
    return
  }

  userId, err := auth.ValidateJWT(token, cfg.jwtsecret)
  if err != nil {
    http.Error(w, err.Error(), http.StatusUnauthorized)
    return
  }

  chirpId, err := uuid.Parse(r.PathValue("id"))
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  deleteChirpParam := database.DeleteChirpParams{
    UserID: userId,
    ID: chirpId,
  }
  chirp, err := cfg.dbQueries.GetChirp(r.Context(), chirpId)
  if err != nil{
    if err == sql.ErrNoRows {
      http.Error(w, err.Error(), http.StatusNotFound)
      return
    }
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  if userId != chirp.UserID {
    http.Error(w, "User is not the Chirp author", http.StatusForbidden)
    return
  }
    
  err = cfg.dbQueries.DeleteChirp(r.Context(), deleteChirpParam)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.WriteHeader(204)
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

func respondWithJson (w http.ResponseWriter, code int, chirps []Chirp) {
  data, err := json.Marshal(chirps)
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

