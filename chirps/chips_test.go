package chirps

import (
  "testing"
)

func TestCleanMessage (t *testing.T) {
  got := sanitizedChirp("I hear Mastodon is better than Chirpy. sharbert I need to migrate")
  want := "I hear Mastodon is better than Chirpy. **** I need to migrate"

  if got != want {
    t.Errorf("got %q, want %q", got, want)
  }
}
