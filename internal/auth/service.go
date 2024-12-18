package auth

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/DanyAdhi/learn-golang/internal/config/redis"
	"github.com/DanyAdhi/learn-golang/internal/users"
	"github.com/DanyAdhi/learn-golang/internal/utils"
)

type Service interface {
	SignUpService(user *UserSignUp) error
	SignIn(data RequestSignIn) (*ResponseSignIn, error)
	RefreshTokenService(refreshToken string) (*ResponseRefreshToken, error)
	SignOutService(userId int, token string) error
}

type service struct {
	repo     Repository
	userRepo users.Repository
}

func NewService(repo Repository, userRepo users.Repository) Service {
	return &service{
		repo:     repo,
		userRepo: userRepo,
	}
}

var ErrWrongEmailOrPassword = errors.New("wrong email or password")
var ErrEmailAlreadyExist = errors.New("email already exist")
var ctx = context.Background()

func (s service) SignUpService(data *UserSignUp) error {
	user, err := s.userRepo.CheckEmailExists(data.Email)
	if err != nil {
		return err
	}
	if user {
		return ErrEmailAlreadyExist
	}

	hashPassword, err := utils.BcryptHashPassword(data.Password)
	if err != nil {
		return err
	}
	data.Password = hashPassword

	err = s.repo.StoreUsersSignUpRepository(data)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) SignIn(data RequestSignIn) (*ResponseSignIn, error) {
	user, err := s.repo.GetUsersByEmail(data.Email)
	if err == sql.ErrNoRows {
		return nil, ErrWrongEmailOrPassword
	}
	if err != nil {
		log.Printf("Failed get user. %v", err)
		return nil, err
	}

	err = utils.CompareHashAndPassword(user.Password, data.Password)
	if err != nil {
		return nil, ErrWrongEmailOrPassword
	}

	payload := utils.PayloadJwt{
		ID:   user.ID,
		Name: user.Name,
	}

	access_token, err := utils.GenerateAccessToken(payload)
	if err != nil {
		return nil, err
	}
	refresh_token, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	err = s.repo.StoreRefreshToken(user.ID, refresh_token)
	if err != nil {
		log.Printf("Failed store refresh token. %v", err)
		return nil, err
	}

	response := &ResponseSignIn{
		Access_token:  access_token,
		Refresh_token: refresh_token,
	}

	return response, nil
}

func (s *service) RefreshTokenService(refreshToken string) (*ResponseRefreshToken, error) {
	err := utils.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// get refresh token in db revoked false
	data, err := s.repo.GetRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// generate new access token
	payload := utils.PayloadJwt{
		ID:   data.User_id,
		Name: data.Name,
	}
	accessToken, err := utils.GenerateAccessToken(payload)
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
