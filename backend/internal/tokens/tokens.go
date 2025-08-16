// File: backend/internal/tokens/tokens.go

package tokens

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

const (
	ScopeAuthentication = "authentication"
	ScopeRefreshToken   = "refresh"
)

var (
	ErrSessionNotFound = errors.New("session not found")
)

// Token is the struct we'll send to the client.
type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    uuid.UUID `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

// SessionData holds the data we store in Redis.
type SessionData struct {
	UserID uuid.UUID `json:"user_id"`
	Roles  []string  `json:"roles"` // Storing roles for RBAC
}

// GenerateToken creates a new token and session in Redis.
func GenerateToken(client *redis.Client, userID uuid.UUID, ttl time.Duration, scope string, roles []string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	// Create 32 random bytes for the token plaintext.
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}
	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	// Hash the token plaintext to use as the Redis key.
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	// Store the session data in Redis.
	session := SessionData{
		UserID: userID,
		Roles:  roles,
	}
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return nil, err
	}

	// Use a context for the Redis operation.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// The key is the hashed token to prevent storing raw tokens.
	// We include the scope to differentiate auth vs. refresh tokens.
	redisKey := fmt.Sprintf("%s:%x", scope, token.Hash)
	err = client.Set(ctx, redisKey, sessionJSON, ttl).Err()
	if err != nil {
		return nil, err
	}

	return token, nil
}

// GetSession retrieves session data from Redis using a token.
func GetSession(client *redis.Client, tokenPlaintext string, scope string) (*SessionData, error) {
	hash := sha256.Sum256([]byte(tokenPlaintext))
	redisKey := fmt.Sprintf("%s:%x", scope, hash[:])

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := client.Get(ctx, redisKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	var session SessionData
	err = json.Unmarshal([]byte(result), &session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// DeleteSession removes a session from Redis. This is our logout functionality.
func DeleteSession(client *redis.Client, tokenPlaintext string, scope string) error {
	hash := sha256.Sum256([]byte(tokenPlaintext))
	redisKey := fmt.Sprintf("%s:%x", scope, hash[:])

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return client.Del(ctx, redisKey).Err()
}
