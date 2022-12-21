// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// HandlerFunc route wrapper for checking JWT token validity
// ----------------------------------------------------------------------------

package auth

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
)

// JWTValidator is a struct that can be used to protect routes
type JWTValidator struct {
	jwksURL string
	scope   string
}

func NewJWTValidator(jwksURL string, scope string) JWTValidator {
	return JWTValidator{
		jwksURL: jwksURL,
		scope:   scope,
	}
}

// Protect can be added around any route handler to protect it and enforce JWT auth
func (v JWTValidator) Protect(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get auth header & bearer scheme
		authHeader := r.Header.Get("Authorization")
		if len(authHeader) == 0 {
			w.WriteHeader(401)
			return
		}

		// Split header into scheme & B64 token	string
		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 {
			w.WriteHeader(401)
			return
		}
		if strings.ToLower(authParts[0]) != "bearer" {
			w.WriteHeader(401)
			return
		}

		// Beyond here we conditionally attempt to fully validate the JWT token

		// Disable check if client id is not set, this is running in demo / unsecured mode
		clientID := os.Getenv("AUTH_CLIENT_ID")
		if clientID == "" {
			log.Println("### Auth: Skipping validation, AUTH_CLIENT_ID is not set")
			next(w, r)
			return
		}

		jwks, err := keyfunc.Get(v.jwksURL, keyfunc.Options{
			RefreshInterval: time.Duration(1) * time.Hour,
		})
		if err != nil {
			log.Printf("### Failed to get the JWKS. Error: %s", err)
			w.WriteHeader(401)
			return
		}

		// Parse the JWT string using the JWKS
		token, err := jwt.Parse(authParts[1], jwks.Keyfunc)
		if err != nil {
			log.Printf("Failed to parse the JWT. Error: %s", err)
			w.WriteHeader(401)
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		// Check the scope includes the app scope
		if !strings.Contains(claims["scp"].(string), v.scope) {
			log.Printf("### Auth: Scope '%s' is missing from token scope '%s'", v.scope, claims["scp"])
			w.WriteHeader(401)
			return
		}

		// Check the token audience is the client id, this might have already been done by jwt.Parse
		if claims["aud"] != clientID {
			log.Printf("### Auth: Token audience '%s' does not match '%s'", claims["aud"], clientID)
			w.WriteHeader(401)
			return
		}

		// Otherwise, we're all good!
		next(w, r)
	}
}
