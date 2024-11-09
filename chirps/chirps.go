package chirps

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
  "fmt"

	"github.com/google/uuid"
)

func WriteChirps(w http.ResponseWriter, r *http.Request) {
  type requestsBody struct {
    Body string `json:"body"`
    UserId uuid.UUID `json:"user_id"`
  }

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
  } else {
    msg, _ := sanitizedChirp(reqBody.Body)
    fmt.Printf("this is the message: %q", msg)
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
