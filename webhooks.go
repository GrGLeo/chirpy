package main

import (
	"encoding/json"
	"net/http"

	"github.com/GrGLeo/chirpy/internal/auth"
	"github.com/google/uuid"
)

type ReqBody struct {
  Event string `json:"event"`
  Data struct {
    UserId uuid.UUID `json:"user_id"`
  } `json:"data"`
}

func (cfg *apiConfig) VerifyPremium (w http.ResponseWriter, r *http.Request) {
  apiKey, _ := auth.GetApiKey(r.Header)

  if apiKey != cfg.apikey {
    http.Error(w, "Not authorize", http.StatusUnauthorized)
  }

  var reqBody ReqBody
  if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
    http.Error(w, "Error decoding body", http.StatusBadRequest)
    return
  }
  if reqBody.Event != "user.upgraded" {
    w.WriteHeader(http.StatusNoContent)
    return
  }
  userID := reqBody.Data.UserId
  success, err := cfg.dbQueries.UpgradeUser(r.Context(), userID)
  if err != nil {
    http.Error(w, "Error upgrading user", http.StatusInternalServerError)
    return
  }
  if success != 1 {
    http.Error(w, "User not found", http.StatusNotFound)
    return
  }
  w.WriteHeader(http.StatusNoContent)
  return
}
