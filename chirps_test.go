package main

import (
  "testing"
)

func TestCleanMessage (t *testing.T) {
  t.Run("testing simple profane word", func(t *testing.T) {
    got, _ := sanitizedChirp("I hear Mastodon is better than Chirpy. sharbert I need to migrate")
    want := "I hear Mastodon is better than Chirpy. **** I need to migrate"

    if got != want {
      t.Errorf("got %q, want %q", got, want)
    }
  })

  t.Run("testing with punctuation", func(t *testing.T) {
    got, _ := sanitizedChirp("I hear Mastodon is better than Chirpy. sharbert! I need to migrate")
    want := "I hear Mastodon is better than Chirpy. ****! I need to migrate"

    if got != want {
      t.Errorf("got %q, want %q", got, want)
    }
  })

  t.Run("testing with lower case", func(t *testing.T) {
    got, _ := sanitizedChirp("I hear Mastodon is better than Chirpy. SHARBERT I need to migrate")
    want := "I hear Mastodon is better than Chirpy. **** I need to migrate"

    if got != want {
      t.Errorf("got %q, want %q", got, want)
    }
  })
}
