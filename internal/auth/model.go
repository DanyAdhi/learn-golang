package auth

import (
	"time"
)

type UserSignUp struct {
	Name     string `json:"name" validate:"required,min=3,alpha"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=50"`
}

type RequestSignIn struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=50"`
}

type ResponseSignIn struct {
	Access_token  string `json:"access_token"`
	Refresh_token string `json:"refresh_token"`
}

type User struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Password  string     `json:"password"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type ReqBodyRefreshToken struct {
	Refresh_token string `json:"refresh_token"`
}

type ResponseRefreshToken struct {
	Access_token string `json:"access_token"`
}

type GetRefreshToken struct {
	User_id int    `json:"id"`
	Name    string `json:"name"`
}
