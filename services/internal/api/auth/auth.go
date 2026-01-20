package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/leroysb/go_kubernetes/internal/cache"
	"github.com/leroysb/go_kubernetes/internal/database/models"
)

var (
	requiredScope     = os.Getenv("REQUIRED_SCOPE")
	hydraAdminURL     = os.Getenv("HYDRA_ADMIN_URL")
	hydraTokenURL     = os.Getenv("HYDRA_TOKEN_URL")
	hydraClientID     = os.Getenv("HYDRA_CLIENT_ID")
	hydraClientSecret = os.Getenv("HYDRA_CLIENT_SECRET")
	hydraScope        = os.Getenv("HYDRA_SCOPE")
	introspectFunc    = introspectToken
)

// SetRequiredScope allows tests to override the required scope dynamically.
func SetRequiredScope(scope string) {
	requiredScope = scope
}

// SetIntrospectFunc allows tests to stub the token introspection behavior.
func SetIntrospectFunc(fn func(string) (*TokenInfo, error)) {
	if fn == nil {
		introspectFunc = introspectToken
		return
	}
	introspectFunc = fn
}

// GetAccessToken requests an access token from ORY Hydra using client credentials
func GetAccessToken() (string, error) {
	if hydraTokenURL == "" || hydraClientID == "" || hydraClientSecret == "" {
		return "", errors.New("hydra client credentials are not configured")
	}

	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	if hydraScope != "" {
		form.Set("scope", hydraScope)
	}

	req, err := http.NewRequest("POST", hydraTokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("creating hydra token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(hydraClientID, hydraClientSecret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("sending hydra token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("hydra token endpoint returned status %d", resp.StatusCode)
		return "", errors.New("token creation failed")
	}

	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", fmt.Errorf("decoding hydra token response: %w", err)
	}

	if tokenResponse.AccessToken == "" {
		return "", errors.New("received empty access token from hydra")
	}

	return tokenResponse.AccessToken, nil
}

// introspectToken sends a request to the Hydra introspection endpoint to validate the access token
func introspectToken(accessToken string) (*TokenInfo, error) {
	if hydraAdminURL == "" {
		return nil, errors.New("hydra admin URL is not configured")
	}

	formData := url.Values{}
	formData.Set("token", accessToken)

	req, err := http.NewRequest("POST", hydraAdminURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("creating introspection request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(hydraClientID, hydraClientSecret)

	introspectResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending introspection request: %w", err)
	}
	defer introspectResp.Body.Close()

	if introspectResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("introspection failed with status code %d", introspectResp.StatusCode)
	}

	var tokenInfo TokenInfo
	if err := json.NewDecoder(introspectResp.Body).Decode(&tokenInfo); err != nil {
		return nil, fmt.Errorf("decoding introspection response: %w", err)
	}

	return &tokenInfo, nil
}

func hasScope(tokenScope string, required string) bool {
	if required == "" {
		return true
	}

	scopes := strings.Split(tokenScope, " ")
	for _, scope := range scopes {
		if scope == required {
			return true
		}
	}
	return false
}

// AuthMiddleware validates access tokens issued by Hydra
func AuthMiddleware(next fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		accessToken := strings.TrimPrefix(authHeader, "Bearer ")
		tokenInfo, err := introspectFunc(accessToken)
		if err != nil || tokenInfo == nil || !tokenInfo.Active {
			log.Printf("token introspection failed: %v", err)
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		if !hasScope(tokenInfo.Scope, requiredScope) {
			log.Printf("token missing required scope %s", requiredScope)
			return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Insufficient scope"})
		}

		sessionCtx := c.UserContext()
		if sessionCtx == nil {
			sessionCtx = context.Background()
		}

		session, err := cache.FetchSession(sessionCtx, accessToken)
		if err != nil {
			log.Printf("failed to fetch session: %v", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
		}

		if session == nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Session expired"})
		}

		customer := &models.Customer{
			Name:  session.Name,
			Phone: session.Phone,
		}
		customer.ID = session.CustomerID

		c.Locals("user", customer)
		c.Locals("token", accessToken)

		return next(c)
	}
}
