package tests

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/leroysb/go_kubernetes/internal/api/auth"
	"github.com/leroysb/go_kubernetes/internal/cache"
	"github.com/leroysb/go_kubernetes/internal/database/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddlewareHydratesUserFromRedis(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	t.Cleanup(func() { mr.Close() })

	prevRedisAddr := os.Getenv("REDIS_ADDR")
	t.Cleanup(func() { os.Setenv("REDIS_ADDR", prevRedisAddr) })
	os.Setenv("REDIS_ADDR", mr.Addr())
	os.Setenv("REDIS_PASSWORD", "")

	cache.Close()
	require.NoError(t, cache.InitRedis())
	t.Cleanup(func() { cache.Close() })

	token := "test-token"
	session := cache.SessionData{CustomerID: 42, Phone: "+123456789", Name: "Test User"}
	require.NoError(t, cache.StoreSession(context.Background(), token, session))

	originalScope := os.Getenv("REQUIRED_SCOPE")
	auth.SetRequiredScope("api")
	t.Cleanup(func() { auth.SetRequiredScope(originalScope) })

	auth.SetIntrospectFunc(func(string) (*auth.TokenInfo, error) {
		return &auth.TokenInfo{Active: true, Scope: "api"}, nil
	})
	t.Cleanup(func() { auth.SetIntrospectFunc(nil) })

	app := fiber.New()
	app.Get("/private", auth.AuthMiddleware(func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(*models.Customer)
		if !ok {
			t.Fatalf("user context missing")
		}
		return c.JSON(user)
	}))

	req := httptest.NewRequest("GET", "/private", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	var payload map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&payload))
	assert.Equal(t, float64(session.CustomerID), payload["ID"])
	assert.Equal(t, session.Name, payload["name"])
	assert.Equal(t, session.Phone, payload["phone"])
}
