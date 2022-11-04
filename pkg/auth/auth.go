// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// HandlerFunc middleware for checking JWT token validity
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

// Fix the JWKS URL to be the one for Azure AD v2
const jwksURL = `https://login.microsoftonline.com/common/discovery/v2.0/keys`

var AppScopeName = "User.Read"

// JWTValidator can be added around any route to protect it
func JWTValidator(next http.HandlerFunc) http.HandlerFunc {

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

		// Just decode the token payload, we don't verify the signature
		// Grab preferred_username claim from token
		username, err := getClaimFromJWT(authParts[1], "preferred_username")
		if err == nil && username != "" {
			r.Header.Set("X-Username", username)
		}

		// Beyond here we fully validate the JWT token

		// Disable check if client id is not set, this is running in demo / unsecured mode
		clientID := os.Getenv("AUTH_CLIENT_ID")
		if clientID == "" {
			log.Println("### Auth: No validation as AUTH_CLIENT_ID is not set")
			next(w, r)
			return
		}

		jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{
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
		if !strings.Contains(claims["scp"].(string), AppScopeName) {
			log.Printf("### Auth: Scope '%s' is missing from token scope '%s'", AppScopeName, claims["scp"])
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
