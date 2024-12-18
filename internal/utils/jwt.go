package utils

import (
	"errors"
	"log"
	"time"

	"github.com/DanyAdhi/learn-golang/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type PayloadJwt struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type JwtDecodeInterface struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(payload PayloadJwt) (string, error) {
	secreatKey := config.AppConfig.JWT_SECRET_ACCESS_TOKEN
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   payload.ID,
		"name": payload.Name,
		"exp":  time.Now().Add(time.Minute * 10).Unix(),
		"iat":  time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(secreatKey))
	if err != nil {
		log.Printf("Failed genrerate access token. %v", err)
		return "", err
	}

	return tokenString, nil
}

func VerifyAccessToken(tokenString string) (*JwtDecodeInterface, error) {
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

func GenerateRefreshToken(id int) (string, error) {
	secretKey := config.AppConfig.JWT_SECRET_REFRESH_TOKEN

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 730).Unix(),
		"iat": time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Printf("Failed generate refresh token. %v", err)
		return "", err
	}
	return tokenString, nil
}

func VerifyRefreshToken(tokenString string) error {
	secretKey := config.AppConfig.JWT_SECRET_REFRESH_TOKEN

	_, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, ErrUnexpectedSigningMethod
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		log.Printf("Error verify refresh tokenssss %v", err)
		return err
	}

	return nil
}
