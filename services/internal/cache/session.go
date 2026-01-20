package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// SessionData represents the stored session payload for an authenticated customer.
type SessionData struct {
	CustomerID uint   `json:"customer_id"`
	Phone      string `json:"phone"`
	Name       string `json:"name"`
}

var (
	ttlOnce sync.Once
	ttl     time.Duration
)

func sessionKey(token string) string {
	return fmt.Sprintf("session:%s", token)
}

func resolveTTL() {
	ttl = time.Hour

	raw := os.Getenv("SESSION_TTL_SECONDS")
	if raw == "" {
		return
	}

	seconds, err := strconv.Atoi(raw)
	if err != nil || seconds <= 0 {
		return
	}

	ttl = time.Duration(seconds) * time.Second
}

// SessionTTL returns the configured TTL for sessions, defaulting to one hour.
func SessionTTL() time.Duration {
	ttlOnce.Do(resolveTTL)
	return ttl
}

// StoreSession persists a session payload keyed by the Hydra access token.
func StoreSession(ctx context.Context, token string, data SessionData) error {
	if Client() == nil {
		return errors.New("redis client is not initialized")
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal session payload: %w", err)
	}

	if err := Client().Set(ctx, sessionKey(token), payload, SessionTTL()).Err(); err != nil {
		return fmt.Errorf("store session: %w", err)
	}

	return nil
}

// FetchSession returns the session payload associated with the access token.
func FetchSession(ctx context.Context, token string) (*SessionData, error) {
	if Client() == nil {
		return nil, errors.New("redis client is not initialized")
	}

	res, err := Client().Get(ctx, sessionKey(token)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("fetch session: %w", err)
	}

	var data SessionData
	if err := json.Unmarshal([]byte(res), &data); err != nil {
		return nil, fmt.Errorf("unmarshal session payload: %w", err)
	}

	return &data, nil
}

// DeleteSession removes the session keyed by the given token.
func DeleteSession(ctx context.Context, token string) error {
	if Client() == nil {
		return errors.New("redis client is not initialized")
	}

	if err := Client().Del(ctx, sessionKey(token)).Err(); err != nil && !errors.Is(err, redis.Nil) {
		return fmt.Errorf("delete session: %w", err)
	}

	return nil
}
