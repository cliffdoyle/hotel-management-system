package passwords

import (
	"golang.org/x/crypto/bcrypt"
)

// Hash generates a bcrypt hash of the password.
func Hash(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

// Matches checks if the provided password matches the stored hash.
func Matches(password string, hash []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		// bcrypt.ErrMismatchedHashAndPassword is the one we want to check for.
		// Other errors might indicate a problem with the hash itself.
		return false, err
	}
	return true, nil
}
