package main

import (
	"encoding/json"
	"net/http"

	"github.com/GrGLeo/chirpy/internal/auth"
	"github.com/GrGLeo/chirpy/internal/database"
)


func (cfg *apiConfig) UpdateUser(w http.ResponseWriter, r *http.Request) {
  token, err := auth.GetBearerToken(r.Header)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  userId, err := auth.ValidateJWT(token, cfg.jwtsecret)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  var updateInfo CreateRequestBody
  if err := json.NewDecoder(r.Body).Decode(&updateInfo); err != nil {
    http.Error(w, "Invalid request body", http.StatusBadRequest)
    return
  }
  
  hashPwd, err := auth.HashPassword(updateInfo.Password)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  updateParams := database.UpdateUserInfoParams{
    HashedPassword: hashPwd,
    Email: updateInfo.Email,
    ID: userId,
  }
  
  updatedUserInfo, err := cfg.dbQueries.UpdateUserInfo(r.Context(), updateParams)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  user := User{
    ID: updatedUserInfo.ID,
    CreatedAt: updatedUserInfo.CreatedAt,
    UpdatedAt: updatedUserInfo.UpdatedAt,
    Email: updatedUserInfo.Email,
  }
  data, err := json.Marshal(user)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusCreated)
  w.Write(data)
}
  
