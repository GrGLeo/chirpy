package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/GrGLeo/chirpy/internal/auth"
	"github.com/GrGLeo/chirpy/internal/database"
)

func (cfg *apiConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
    type RequestBody struct {
      Password string `json:"password"`
      Email string `json:"email"` 
  }
  
  type User struct {
    ID        uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Email     string    `json:"email"`
  }

  var reqBody RequestBody
  if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
    http.Error(w, "Invalid request body", http.StatusBadRequest)
    return
  }
  
  hashedPassword, err := auth.HashPassword(reqBody.Password)
  if err != nil {
    http.Error(w, "Internal server error", http.StatusInternalServerError)
  }

  userParam := database.CreateUserParams{
    Email: reqBody.Email,
    HashedPassword: hashedPassword,
  }

  newUser, err := cfg.dbQueries.CreateUser(r.Context(), userParam)
  user := User{
    ID: newUser.ID,
    CreatedAt: newUser.CreatedAt,
    UpdatedAt: newUser.UpdatedAt,
    Email: newUser.Email,
  }
  
  if err != nil {
    http.Error(w, "Error while creating user", http.StatusConflict)
    return
  }

  data, err := json.Marshal(user)
  if err != nil {
    log.Printf("Error marshaling json: %s", err)
    w.WriteHeader(500)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusCreated)
  w.Write(data)
}
