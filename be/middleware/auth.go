package middleware

import (
	"crypto/rsa"
	"net/http"
	"strings"
	"be/constant"
	"be/response"
	"github.com/dgrijalva/jwt-go/v4"
)

var parsedPublicKey *rsa.PublicKey

func initialize() {
    if constant.APPCONFIG == nil {
        log.Fatal("Configuration is not initialized.")
    }
    if constant.APPCONFIG.PublicKey == "" {
        log.Fatal("Public key is missing in configuration")
    }
    // Parse the public key from the configuration
    var err error
    parsedPublicKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(constant.APPCONFIG.PublicKey))
    if err != nil {
        log.Fatal("Error parsing public key: ", err)
    }
}


func VerifyToken(next http.Handler) http.Handler {
	initialize()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
            log.Error("Authorization header is missing");
			response.PutInvalidJwtResponse("Authorization header is missing", w)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			log.Error("Bearer token is missing");
			response.PutInvalidJwtResponse("Bearer token is missing", w)
			return
		}

		// Prepare to parse the token
		claims := &jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return parsedPublicKey, nil
		})

		// Handle token parsing errors
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				log.Error("Invalid token signature");
				response.PutInvalidJwtResponse("Invalid token signature", w)
				return
			}
			log.Error("Invalid token");
			response.PutInvalidJwtResponse("Invalid token", w)
			return
		}

		// Check if token is valid
		if !token.Valid {
			log.Error("Invalid token");
			response.PutInvalidJwtResponse("Invalid token", w)
			return
		}

		// Token is valid; proceed with the request
		next.ServeHTTP(w, r)
	})
}
