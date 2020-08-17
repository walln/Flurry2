package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/dgrijalva/jwt-go"
	"github.com/walln/flurry2/flurry/global"
)

func AuthenticateWithFirebase(w http.ResponseWriter, r *http.Request) bool {

	logger := log.WithFields(log.Fields{
		"Proxy Route":    r.URL.RequestURI(),
		"Client IP":      r.RemoteAddr,
		"Request Method": r.Method,
	})

	authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
	proceed := true
	if len(authHeader) != 2 {
		proceed = false
	} else {
		client, err := global.GetFirebaseApp().Auth(r.Context())
		if err != nil {
			logger.Fatalf("Error getting Auth client: %v\n", err)
			proceed = false
		}
		_, err = client.VerifyIDToken(r.Context(), authHeader[1])
		if err != nil {
			logger.Printf("Error verifying ID token: %v\n", err)
			proceed = false

		}
	}
	return proceed
}

func AuthenticateWithJWT(w http.ResponseWriter, r *http.Request) bool {

	logger := log.WithFields(log.Fields{
		"Proxy Route":    r.URL.RequestURI(),
		"Client IP":      r.RemoteAddr,
		"Request Method": r.Method,
	})

	authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
	proceed := true
	if len(authHeader) != 2 {
		logger.Info("Request contains malformed authorization token. Rejecting forward request.")
		proceed = false

	} else {
		SECRETKEY := os.Getenv("JWT_SECRET")
		if SECRETKEY == "" {
			panic("JWT Secret not supplied!")
		}
		jwtToken := authHeader[1]
		token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
			if global.GetSigningMethod() == "HMAC" {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
			} else if global.GetSigningMethod() == "ECDSA" {
				if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
			} else if global.GetSigningMethod() == "RSA" {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
			} else if global.GetSigningMethod() == "RSAPSS" {
				if _, ok := token.Method.(*jwt.SigningMethodRSAPSS); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
			} else {
				panic("No signing method specified!")
			}
			return []byte(SECRETKEY), nil
		})

		if !token.Valid {
			logger.Debug("Request contains Invalid token: " + err.Error())
			proceed = false
		}
	}
	return proceed
}
