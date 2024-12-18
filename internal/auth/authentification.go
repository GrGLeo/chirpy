package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pwd string) (string, error) {
  encryptPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), 1)
  if err != nil {
    fmt.Println("Error while encrypting password")
  }
  return string(encryptPwd), nil
}


func CheckPasswordHash(pwd, hash string) error {
  err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
  if err != nil {
    return err
  }
  return nil
}

func MakeJWT(userId uuid.UUID, tokenSecret string) (string, error) {
  now := time.Now()
  expiresIn := time.Duration(1) * time.Hour
  
  claims := &jwt.RegisteredClaims{
    Issuer: "chirpy",
    IssuedAt: jwt.NewNumericDate(now),
    ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
    Subject: userId.String(),
  }
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
  type CustomClaims struct {
    jwt.RegisteredClaims
  }

  token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
    return []byte(tokenSecret), nil
  })
  if err != nil {
    fmt.Println("glop")
    return uuid.Nil, err
  }

  if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
    userID, err := uuid.Parse(claims.Subject)
    if err != nil {
      return uuid.Nil, err
    }
    return userID, nil
  }
  return uuid.Nil, fmt.Errorf("invalid token")
}

func GetBearerToken(header http.Header) (string, error) {
  authorization := header.Get("Authorization")
  if authorization == "" {
    return authorization, errors.New("No authorization found")
  }
  bearerToken := strings.Split(authorization, " ")
  return bearerToken[1], nil
}

func MakeRefreshToken() (string, error) {
  randomData := make([]byte, 32)
  _, err := rand.Read(randomData)
  if err != nil {
    return "", err
    }
  return hex.EncodeToString(randomData), nil
}

func GetApiKey(header http.Header) (string, error) {
  authorization := header.Get("Authorization")
  if authorization == "" {
    return authorization, errors.New("No authorization found")
  }
  apiKey := strings.Split(authorization, " ")
  return apiKey[1], nil
}
