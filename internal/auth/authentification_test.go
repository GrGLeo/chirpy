package auth

import (
  "testing"
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

