package auth

import "time"

type ClientRequest struct {
	ClientName        string   `json:"client_name"`
	ClientSecret      string   `json:"client_secret"`
	GrantTypes        []string `json:"grant_types"`
	Scope             string   `json:"scope"`
	TokenEndpointAuth string   `json:"token_endpoint_auth_method"`
}

type ClientResponse struct {
	ClientID              string   `json:"client_id"`
	ClientName            string   `json:"client_name"`
	ClientSecret          string   `json:"client_secret"`
	RedirectUris          any      `json:"redirect_uris"`
	GrantTypes            []string `json:"grant_types"`
	ResponseTypes         any      `json:"response_types"`
	Scope                 string   `json:"scope"`
	Audience              []any    `json:"audience"`
	Owner                 string   `json:"owner"`
	PolicyURI             string   `json:"policy_uri"`
	AllowedCorsOrigins    []any    `json:"allowed_cors_origins"`
	TosURI                string   `json:"tos_uri"`
	ClientURI             string   `json:"client_uri"`
	LogoURI               string   `json:"logo_uri"`
	Contacts              any      `json:"contacts"`
	ClientSecretExpiresAt int      `json:"client_secret_expires_at"`
	SubjectType           string   `json:"subject_type"`
	Jwks                  struct {
	} `json:"jwks"`
	TokenEndpointAuthMethod   string    `json:"token_endpoint_auth_method"`
	UserinfoSignedResponseAlg string    `json:"userinfo_signed_response_alg"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
	Metadata                  struct {
	} `json:"metadata"`
	RegistrationAccessToken                    string `json:"registration_access_token"`
	RegistrationClientURI                      string `json:"registration_client_uri"`
	SkipConsent                                bool   `json:"skip_consent"`
	SkipLogoutConsent                          any    `json:"skip_logout_consent"`
	AuthorizationCodeGrantAccessTokenLifespan  any    `json:"authorization_code_grant_access_token_lifespan"`
	AuthorizationCodeGrantIDTokenLifespan      any    `json:"authorization_code_grant_id_token_lifespan"`
	AuthorizationCodeGrantRefreshTokenLifespan any    `json:"authorization_code_grant_refresh_token_lifespan"`
	ClientCredentialsGrantAccessTokenLifespan  any    `json:"client_credentials_grant_access_token_lifespan"`
	ImplicitGrantAccessTokenLifespan           any    `json:"implicit_grant_access_token_lifespan"`
	ImplicitGrantIDTokenLifespan               any    `json:"implicit_grant_id_token_lifespan"`
	JwtBearerGrantAccessTokenLifespan          any    `json:"jwt_bearer_grant_access_token_lifespan"`
	RefreshTokenGrantIDTokenLifespan           any    `json:"refresh_token_grant_id_token_lifespan"`
	RefreshTokenGrantAccessTokenLifespan       any    `json:"refresh_token_grant_access_token_lifespan"`
	RefreshTokenGrantRefreshTokenLifespan      any    `json:"refresh_token_grant_refresh_token_lifespan"`
}

type TokenRequest struct {
	ClientID  string `json:"client_id"`
	GrantType string `json:"grant_type"`
	// ClientSecret string `json:"client_secret"`
	// Scope        string `json:"scope"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

type TokenInfo struct {
	Active    bool   `json:"active"`
	Scope     string `json:"scope"`
	ClientID  string `json:"client_id"`
	Sub       string `json:"sub"`
	Exp       int    `json:"exp"`
	Iat       int    `json:"iat"`
	Nbf       int    `json:"nbf"`
	Aud       []any  `json:"aud"`
	Iss       string `json:"iss"`
	TokenType string `json:"token_type"`
	TokenUse  string `json:"token_use"`
}
