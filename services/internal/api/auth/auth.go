package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var requiredScope = os.Getenv("requiredScope")
var hydraAdminUrl = os.Getenv("hydraAdminUrl")

// HydraClientResponse communicates with the Hydra admin API to create a new OAuth2 client
func GetAccessToken() (string, error) {
	var clientURL = os.Getenv("HYDRA_CLIENT_URL")
	var tokenURL = os.Getenv("HYDRA_TOKEN_URL")
	var scope = os.Getenv("HYDRA_SCOPE")
	var clientName = os.Getenv("HYDRA_CLIENT_NAME")
	var clientSecret = os.Getenv("HYDRA_CLIENT_SECRET")

	client := &http.Client{}

	var clientRequest = ClientRequest{
		ClientName:        clientName,
		ClientSecret:      clientSecret,
		GrantTypes:        []string{"authorization_code", "refresh_token"},
		Scope:             scope,
		TokenEndpointAuth: "none",
	}

	clientRequestBody, err := json.Marshal(clientRequest)
	if err != nil {
		fmt.Println("Error marshalling client id request:", err)
		return "", err
	}

	// create the HTTP request
	clientReq, err := http.NewRequest("POST", clientURL, bytes.NewBuffer(clientRequestBody))
	if err != nil {
		fmt.Println("Error creating client id request:", err)
		return "", err
	}

	// set the request headers
	clientReq.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Accept", "*/*")
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	// req.Header.Set("Connection", "keep-alive")

	// send the request
	clientResp, err := client.Do(clientReq)
	if err != nil {
		fmt.Println("Error sending client id request:", err)
		return "", err
	}
	defer clientResp.Body.Close()

	if clientResp.StatusCode != http.StatusOK {
		return "", errors.New("client id creation failed. unexpected status code")
	}

	var clientResponse ClientResponse
	if err := json.NewDecoder(clientResp.Body).Decode(&clientResponse); err != nil {
		fmt.Println("Error decoding client id response:", err)
		return "", err
	}

	// TO-DO: Delete
	fmt.Println("Client ID:", clientResponse.ClientID)

	var tokenRequest = TokenRequest{
		ClientID:  clientResponse.ClientID,
		GrantType: "client_credentials",
		// ClientSecret: clientSecret,
		// Scope:        scope,
	}

	var tokenRequestBody []byte
	tokenRequestBody, err = json.Marshal(tokenRequest)
	if err != nil {
		fmt.Println("Error marshalling token request:", err)
		return "", err
	}

	// create the HTTP request for token creation
	var tokenReq *http.Request
	tokenReq, err = http.NewRequest("POST", tokenURL, bytes.NewBuffer(tokenRequestBody))
	if err != nil {
		fmt.Println("Error creating token request:", err)
		return "", err
	}

	// set the request headers
	tokenReq.Header.Set("Content-Type", "application/json")
	// tokenReq.Header.Set("Accept", "*/*")
	// tokenReq.Header.Set("Accept-Encoding", "gzip, deflate, br")
	// tokenReq.Header.Set("Connection", "keep-alive")

	var tokenResp *http.Response
	tokenResp, err = client.Do(tokenReq)
	if err != nil {
		fmt.Println("Error sending token request:", err)
		return "", err
	}
	defer tokenResp.Body.Close()

	if tokenResp.StatusCode != http.StatusOK {
		return "", errors.New("token creation failed")
	}

	var tokenResponse TokenResponse
	if err := json.NewDecoder(tokenResp.Body).Decode(&tokenResponse); err != nil {
		fmt.Println("Error decoding token response:", err)
		return "", err
	}

	// TO-DO: Delete
	fmt.Println("Access Token:", tokenResponse.AccessToken)

	return tokenResponse.AccessToken, nil
}

// introspectToken sends a request to the Hydra introspection endpoint to validate the access token
func introspectToken(accessToken string) (*TokenInfo, error) {
	// Prepare the form data
	formData := url.Values{}
	formData.Set("token", accessToken)
	formData.Set("scope", requiredScope)

	// Create the HTTP request
	req, err := http.NewRequest("POST", hydraAdminUrl, strings.NewReader(formData.Encode()))
	if err != nil {
		fmt.Println("Error creating introspection request:", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	// Send the request
	introspectResp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending introspection request:", err)
		return nil, err
	}
	defer introspectResp.Body.Close()

	// Check the response status code
	if introspectResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("introspection failed with status code %d", introspectResp.StatusCode)
	}

	// Parse introspection response
	var tokenInfo TokenInfo
	if err := json.NewDecoder(introspectResp.Body).Decode(&tokenInfo); err != nil {
		return nil, err
	}

	return &tokenInfo, nil
}

// checkTokenActive checks if the introspected token is active
func hasScope(tokenScope string, requiredScope string) bool {
	scopes := strings.Split(tokenScope, " ")
	for _, scope := range scopes {
		if scope == requiredScope {
			return true
		}
	}
	return false
}

// AuthMiddleware is a middleware function to validate access token using Hydra introspection endpoint
func AuthMiddleware(next fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the access token from the request headers
		authHeader := c.Get("Authorization")

		// Check if Authorization header is missing or does not start with "Bearer "
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		// Extract the token by stripping the "Bearer " prefix
		accessToken := strings.TrimPrefix(authHeader, "Bearer ")

		// Introspect the token
		tokenInfo, err := introspectToken(accessToken)
		if err != nil || !tokenInfo.Active {
			fmt.Println("Error during token introspection:", err)
			// return err
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		// Check if the token is active
		if !hasScope(tokenInfo.Scope, requiredScope) {
			// If token is not active or not present, return unauthorized error
			fmt.Printf("Insufficient scope: %v\n", err)
			// return err
			return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Insufficient scope"})
		}

		// Proceed to the next middleware or handler if token is valid
		return c.Next()
	}
}
