package shopify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"net/http"
	"net/url"

	"github.com/golang-jwt/jwt/v5"
)

const JWT_PERMITTED_CLOCK_TOLERANCE = 20 * time.Second // Clock tolerance of 15 seconds

var (
	APIKey       string
	APISecretKey string
)

// init function to initialize global variables
func init() {
	APIKey = os.Getenv("SHOPIFY_API_KEY")
	APISecretKey = os.Getenv("SHOPIFY_API_SECRET")

	// Ensure APIKey and APISecretKey are provided
	if APIKey == "" || APISecretKey == "" {
		log.Fatalf("SHOPIFY_API_SECRET or SHOPIFY_API_SECRET environment variables are not set")
	}
}

type ShopifyClaims struct {
	Dest                 string `json:"dest,omitempty"` // Custom destination field
	jwt.RegisteredClaims        // Embedding standard claims (exp, aud, etc.)
}

type DecodeSessionTokenResponse struct {
	Claims          ShopifyClaims
	MyshopifyDomain string
}

// DecodeSessionToken decodes and verifies the JWT session token.
func DecodeSessionToken(token string) (*DecodeSessionTokenResponse, error) {
	// Parse and validate the token
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(APISecretKey), nil
	}

	tokenObj, err := jwt.ParseWithClaims(token, &ShopifyClaims{}, keyFunc, jwt.WithLeeway(JWT_PERMITTED_CLOCK_TOLERANCE))
	if err != nil {
		return nil, fmt.Errorf("failed to parse session token '%s': %v", token, err)
	}

	// Extract the payload
	claims, ok := tokenObj.Claims.(*ShopifyClaims)
	if !ok || !tokenObj.Valid {
		return nil, fmt.Errorf("invalid session token")
	}

	// Check if the audience matches the API key
	if len(claims.Audience) == 0 || claims.Audience[0] != APIKey {
		return nil, fmt.Errorf("session token had invalid API key")
	}

	parsedURL, err := url.Parse(claims.Dest)
	if err != nil {
		log.Printf("Error parsing myshopify URL: %v", err)
		return nil, errors.New("invalid myshopify URL")
	}
	resp := &DecodeSessionTokenResponse{
		Claims:          *claims,
		MyshopifyDomain: parsedURL.Hostname(),
	}

	return resp, nil
}

// ExchangeTokenResponse represents the structure of the response from Shopify API.
type ExchangeTokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
}

// GetAccessTokenFromShopify exchanges the session token for an offline access token and returns the full ExchangeTokenResponse.
func GetAccessTokenFromShopify(shop string, sessionToken string) (*ExchangeTokenResponse, error) {
	// Define the URL
	url := fmt.Sprintf("https://%s/admin/oauth/access_token", shop)

	// Define the request body with the variable values
	requestBody := map[string]string{
		"client_id":            APIKey,       // Use APIKey for client_id
		"client_secret":        APISecretKey, // Use APISecretKey for client_secret
		"grant_type":           "urn:ietf:params:oauth:grant-type:token-exchange",
		"subject_token":        sessionToken,
		"subject_token_type":   "urn:ietf:params:oauth:token-type:id_token",
		"requested_token_type": "urn:shopify:params:oauth:token-type:offline-access-token",
	}

	// Marshal the payload into JSON
	jsonPayload, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %v", err)
	}

	// Create the POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Send the request using http.Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read and handle the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var tokenResponse ExchangeTokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	// Return the full token response
	return &tokenResponse, nil
}
