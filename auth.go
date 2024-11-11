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

type CreateRequestBody struct {
  Password string `json:"password"`
  Email string `json:"email"` 
}

type LoginRequestBody struct {
  Password string `json:"password"`
  Email string `json:"email"` 
  ExpiresIn int16 `json:"expires_in_seconds"`
}

type User struct {
  ID        uuid.UUID `json:"id"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
  Email     string    `json:"email"`
}

type UserLogged struct {
  ID        uuid.UUID `json:"id"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
  Email     string    `json:"email"`
  Token     string    `json:"token"`
}

func (cfg *apiConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
  var reqBody CreateRequestBody
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

func (cfg *apiConfig) UserLogin(w http.ResponseWriter, r *http.Request) {
  var reqBody LoginRequestBody
  if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
    http.Error(w, "Invaled request body", http.StatusBadRequest)
    return
  }

  userInfo, err := cfg.dbQueries.GetHashedPassword(r.Context(), reqBody.Email)
  if err != nil {
    http.Error(w, "Users not found", http.StatusBadRequest)
    return
  }

  err = auth.CheckPasswordHash(reqBody.Password, userInfo.HashedPassword)
  if err != nil {
    http.Error(w, "Incorrect password", http.StatusUnauthorized)
    return
  }

  // Create JWT token
  expiresIn := reqBody.ExpiresIn
  var expires time.Duration
  if expiresIn == 0 {
    expires = time.Duration(1) * time.Hour
  } else if expiresIn > 3600 {
    expires = time.Duration(1) * time.Hour
  } else {
    expires = time.Duration(expiresIn) * time.Second
  }
  token, err := auth.MakeJWT(userInfo.ID, cfg.jwtsecret, expires)
  if err != nil {
    http.Error(w, "Internal server error", http.StatusInternalServerError)
  }


  user := UserLogged{
    ID: userInfo.ID,
    CreatedAt: userInfo.CreatedAt,
    UpdatedAt: userInfo.UpdatedAt,
    Email: userInfo.Email,
    Token: token,
  }
  
  data, err := json.Marshal(user)
  if err != nil {
    http.Error(w, "Error marshalling json", http.StatusInternalServerError)
    return
  }

  
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  w.Write(data)
}
