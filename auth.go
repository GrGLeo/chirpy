package main

import (
	"database/sql"
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
  IsChirpRed bool     `json:"is_chirpy_red"`
}

type UserLogged struct {
  ID            uuid.UUID `json:"id"`
  CreatedAt     time.Time `json:"created_at"`
  UpdatedAt     time.Time `json:"updated_at"`
  Email         string    `json:"email"`
  Token         string    `json:"token"`
  RefreshToken  string    `json:"refresh_token"`
  IsChirpRed    bool      `json:"is_chirpy_red"`
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
    IsChirpRed: false,
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
  token, err := auth.MakeJWT(userInfo.ID, cfg.jwtsecret)
  if err != nil {
    http.Error(w, "Internal server error", http.StatusInternalServerError)
  }

  // Create RefreshToken
  refreshtoken, err := auth.MakeRefreshToken()
  if err != nil {
    http.Error(w, "Internal server error", http.StatusInternalServerError)
  }
  
  durationHour := 60 * 24
  expiresAt := time.Now().Add(time.Duration(durationHour) * time.Hour)
  refreshTokenParams := database.WriteRefreshTokenParams{
    Token: refreshtoken,
    UserID: userInfo.ID,
    ExpiresAt: expiresAt,
  }

  err = cfg.dbQueries.WriteRefreshToken(r.Context(), refreshTokenParams)
  if err != nil {
    http.Error(w, "Internal server error", http.StatusInternalServerError)
  }

  user := UserLogged{
    ID: userInfo.ID,
    CreatedAt: userInfo.CreatedAt,
    UpdatedAt: userInfo.UpdatedAt,
    Email: userInfo.Email,
    Token: token,
    RefreshToken: refreshtoken,
    IsChirpRed: userInfo.IsChirpyRed,
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

func (cfg *apiConfig) RefreshToken(w http.ResponseWriter, r *http.Request) {
  refreshToken, err := auth.GetBearerToken(r.Header)
  if err != nil {
    http.Error(w, err.Error(), http.StatusUnauthorized)
    return
  }

  refreshTokenRow, err := cfg.dbQueries.GetRefreshToken(r.Context(), refreshToken)
  if err != nil {
    if err == sql.ErrNoRows {
      http.Error(w, "Token not found", http.StatusUnauthorized)
      return
    }
    http.Error(w, "Internal server error", http.StatusInternalServerError)
  }
 // Check if token is still valid 
  now := time.Now()
  if now.After(refreshTokenRow.ExpiresAt) {
    http.Error(w, "Token past valid date", http.StatusUnauthorized)
  }
  if refreshTokenRow.RevokedAt.Valid {
      http.Error(w, "Token is revoked", http.StatusUnauthorized)
      return
  }
  token, _:= auth.MakeJWT(refreshTokenRow.UserID, cfg.jwtsecret)
  
  data, err := json.Marshal(map[string]string{"token": token})
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  w.Write(data)
}

func (cfg *apiConfig) RevokeToken (w http.ResponseWriter, r *http.Request) {
  refreshToken, err := auth.GetBearerToken(r.Header)
  if err != nil {
    http.Error(w, err.Error(), http.StatusUnauthorized)
    return
  }

  revokeTokenParams := database.RevokeTokenParams{
    UpdatedAt: time.Now(),
    Token: refreshToken,
  }

  err = cfg.dbQueries.RevokeToken(r.Context(), revokeTokenParams)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.WriteHeader(http.StatusNoContent)
}
