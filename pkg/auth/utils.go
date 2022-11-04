package auth

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"strings"
)

func getClaimFromJWT(jwt string, claimName string) (string, error) {
	jwtParts := strings.Split(jwt, ".")

	// decode base64 token
	tokenBytes, err := base64.RawURLEncoding.DecodeString(jwtParts[1])
	if err != nil {
		log.Println("### Auth: Error decoding token", err)
		return "", err
	}

	// decode token JSON
	var tokenJSON map[string]interface{}
	err = json.Unmarshal(tokenBytes, &tokenJSON)
	if err != nil {
		return "", err
	}

	// get claim
	claim, ok := tokenJSON[claimName]
	if !ok {
		return "", err
	}

	return claim.(string), nil
}
