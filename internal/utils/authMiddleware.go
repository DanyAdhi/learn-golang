package utils

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/DanyAdhi/learn-golang/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

var ErrUnexpectedSigningMethod = errors.New("unexpected signing method")

type JwtDecodeInterface struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	jwt.RegisteredClaims
}

type contextKey string

const userKey contextKey = "user"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			token := getBearerToken(r)
			if token == "" {
				ResponseError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			decode, err := verifyAccessToken(token)
			if err == jwt.ErrTokenExpired {
				ResponseError(w, http.StatusUnauthorized, "Token expired")
				return
			}
			if err != nil {
				ResponseError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			ctx := context.WithValue(r.Context(), userKey, decode)

			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}

func getBearerToken(r *http.Request) string {
	authorization := r.Header.Get("Authorization")
	if authorization == "" || !strings.Contains(authorization, "Bearer") {
		return ""
	}
	token := strings.Split(authorization, " ")[1]
	return token
}

func verifyAccessToken(tokenString string) (*JwtDecodeInterface, error) {
	claims := &JwtDecodeInterface{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, ErrUnexpectedSigningMethod
		}
		return []byte(config.AppConfig.JWT_SECRET_ACCESS_TOKEN), nil
	})

	if errors.Is(err, jwt.ErrTokenExpired) {
		return nil, jwt.ErrTokenExpired
	}

	if err != nil {
		log.Printf("Error verifying token: %v", err)
		return nil, err
	}

	return claims, nil
}
