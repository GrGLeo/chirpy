package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestPassword (t *testing.T) {
  t.Run("testing password hashing", func(t *testing.T) {
    pwd := "hello"
    hashedPwd, err := HashPassword(pwd)
    if err != nil {
      t.Fatal(err)
    }
    
    if pwd == hashedPwd {
      t.Error("Password shouldn't be equal to the hashed password.")
    }

    err = CheckPasswordHash(pwd, hashedPwd)
    if err != nil {
      t.Error("Password should match")
    }
    
    err = CheckPasswordHash("wrongpwd", hashedPwd)
    if err == nil {
      t.Error("Password shouldn't match")
    }
  })
}

func TestJWT (t *testing.T) {
  t.Run("testing JWT creation", func(t *testing.T) {
    userId := uuid.New()
    expiresIn := time.Duration(1) * time.Hour
    token, err := MakeJWT(userId, "secret", expiresIn)
    if err != nil {
      t.Errorf("Error while creating the token: %q", err)
      return
    }

    retrievedUserId, err := ValidateJWT(token, "secret")
    if err != nil {
      t.Errorf("Error with JWT validation, error: %q", err)
    }
    if retrievedUserId.String() != userId.String() {
      t.Error("UUID dit not match")
    }
  })
}

func TestGetBearerToken (t *testing.T) {
  header := http.Header{"Authorization": []string{"Bearer testing"}}
  token, err := GetBearerToken(header)
 
  if err != nil {
    t.Errorf("Error with authorization: %q", err)
  }

  if token != "testing" {
    t.Errorf("got: %q, want %q", token, "testing")
  }
}
