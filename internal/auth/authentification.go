package auth

import (
	"fmt"
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

func MakeJWT(userId uuid.UUID, tokenSecret string, expiredIn time.Duration) (string, error) {
  time := time.Now()
  
  claims := &jwt.RegisteredClaims{
    Issuer: "chirpy",
    IssuedAt: jwt.NewNumericDate(time),
    ExpiresAt: jwt.NewNumericDate(time.Add(expiredIn)),
    Subject: userId.String(),
  }
  
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  return token.SignedString(tokenSecret)
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
  type CustomClaims struct {
    jwt.RegisteredClaims
  }

  token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
    return []byte(tokenString), nil
  })
  
  if err != nil {
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
