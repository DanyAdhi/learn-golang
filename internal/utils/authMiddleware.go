package utils

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/DanyAdhi/learn-golang/internal/config/redis"
	"github.com/golang-jwt/jwt/v5"
)

var ErrUnexpectedSigningMethod = errors.New("unexpected signing method")

type contextKey string

var UserKey contextKey = "user"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			token := getBearerToken(r)
			if token == "" {
				ResponseError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			decode, err := VerifyAccessToken(token)
			if err == jwt.ErrTokenExpired {
				ResponseError(w, http.StatusUnauthorized, "Token expired")
				return
			}
			if err != nil {
				ResponseError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			checkExpired := checkTokenExpiredOnRedis(token)
			if checkExpired {
				log.Print("Token already logout")
				ResponseError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			ctx := context.WithValue(r.Context(), UserKey, decode)

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

func checkTokenExpiredOnRedis(token string) bool {
	redis := redis.Connect()
	checkExpired := redis.Get(context.Background(), "expired-"+token).Val()
	return checkExpired == "1"
}
