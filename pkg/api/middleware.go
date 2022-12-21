// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Middleware and optional extra routes available to any API
// ----------------------------------------------------------------------------

package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/cors"
)

// Add JWT username to context
// Tries to get the username from the JWT token, but ignores any errors or missing token
func (b *Base) JWTUsernameMiddleware(claim string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if len(authHeader) == 0 {
				next.ServeHTTP(w, r)
				return
			}
			authParts := strings.Split(authHeader, " ")
			if len(authParts) != 2 {
				next.ServeHTTP(w, r)
				return
			}
			if strings.ToLower(authParts[0]) != "bearer" {
				next.ServeHTTP(w, r)
				return
			}

			username, err := getClaimFromJWT(authParts[1], claim)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), "username", username)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

// SimpleCORSMiddleware adds permissive CORS headers to all responses
func (b *Base) SimpleCORSMiddleware(next http.Handler) http.Handler {
	log.Printf("### 🎭 API: configured simple CORS")

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cors.Handler(next).ServeHTTP(w, r)
	})
}

// getClaimFromJWT is a helper to return a claim from a JWT
// It decodes the raw JWT, parses the JSON and returns the claim
func getClaimFromJWT(jwtRaw string, claimName string) (string, error) {
	jwtParts := strings.Split(jwtRaw, ".")

	// Decode base64 main part of the token
	tokenBytes, err := base64.RawURLEncoding.DecodeString(jwtParts[1])
	if err != nil {
		log.Println("### Auth: Error in base64 decoding token", err)
		return "", err
	}

	// Parse token JSON
	var tokenJSON map[string]interface{}
	err = json.Unmarshal(tokenBytes, &tokenJSON)
	if err != nil {
		log.Println("### Auth: Error in JSON parsing token", err)
		return "", err
	}

	// Get the claim
	claim, ok := tokenJSON[claimName]
	if !ok {
		log.Println("### Auth: Claim not found in token", err)
		return "", err
	}

	return claim.(string), nil
}
