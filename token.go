package main

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

func getToken(issuerID string, keyID string) *jwt.Token {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss": issuerID,
		"exp": time.Now().Add(time.Minute * 5).Unix(),
		"aud": "appstoreconnect-v1",
	})

	token.Header = map[string]interface{}{
		"kid": keyID,
		"alg": "ES256",
		"typ": "JWT",
	}

	return token
}
