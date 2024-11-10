package auth

import (
	"fmt"

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
