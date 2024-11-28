package auth

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/DanyAdhi/learn-golang/internal/config"
	"github.com/DanyAdhi/learn-golang/internal/config/redis"
	"github.com/DanyAdhi/learn-golang/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Login(data RequestLogin) (*ResponseLogin, error)
	RefreshTokenService(refreshToken string) (*ResponseRefreshToken, error)
	SignOutService(userId int, token string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

var ErrWrongEmailOrPassword = errors.New("wrong email or password")
var ctx = context.Background()

func (s *service) Login(data RequestLogin) (*ResponseLogin, error) {
	user, err := s.repo.GetUsersByEmail(data.Email)
	if err == sql.ErrNoRows {
		return nil, ErrWrongEmailOrPassword
	}
	if err != nil {
		log.Printf("Failed get user. %v", err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		return nil, ErrWrongEmailOrPassword
	}

	payload := PayloadJwt{
		ID:   user.ID,
		Name: user.Name,
	}

	access_token, err := generateAccessToken(payload)
	if err != nil {
		return nil, err
	}
	refresh_token, err := generateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	err = s.repo.StoreRefreshToken(user.ID, refresh_token)
	if err != nil {
		log.Printf("Failed store refresh token. %v", err)
		return nil, err
	}

	response := &ResponseLogin{
		Access_token:  access_token,
		Refresh_token: refresh_token,
	}

	return response, nil
}

func (s *service) RefreshTokenService(refreshToken string) (*ResponseRefreshToken, error) {
	err := verifyRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// get refresh token in db revoked false
	data, err := s.repo.GetRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// generate new access token
	payload := PayloadJwt{
		ID:   data.User_id,
		Name: data.Name,
	}
	accessToken, err := generateAccessToken(payload)
	if err != nil {
		return nil, err
	}

	result := &ResponseRefreshToken{
		Access_token: accessToken,
	}

	return result, nil
}

func (s *service) SignOutService(userId int, token string) error {
	err := s.repo.RevokeToken(userId)
	if err != nil {
		log.Printf("Error revoke token. %v", err)
		return err
	}

	redis := redis.Connect()
	redis.Set(ctx, "expired-"+token, true, time.Minute*10)

	return nil
}

func generateAccessToken(payload PayloadJwt) (string, error) {
	secretKey := config.AppConfig.JWT_SECRET_ACCESS_TOKEN
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   payload.ID,
		"name": payload.Name,
		"exp":  time.Now().Add(time.Minute * 10).Unix(),
		"iat":  time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Printf("Failed genrerate access token. %v", err)
		return "", err
	}

	return tokenString, nil
}

func generateRefreshToken(id int) (string, error) {
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

func verifyRefreshToken(tokenString string) error {
	secretKey := config.AppConfig.JWT_SECRET_REFRESH_TOKEN

	_, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, utils.ErrUnexpectedSigningMethod
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		log.Printf("Error verify refresh tokenssss %v", err)
		return err
	}

	return nil
}
