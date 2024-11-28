package auth

import "time"

type RequestLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResponseLogin struct {
	Access_token  string `json:"access_token"`
	Refresh_token string `json:"refresh_token"`
}

type User struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Password  string     `json:"password"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type PayloadJwt struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ResponseRefreshToken struct {
	Access_token string `json:"access_token"`
}

type GetRefreshToken struct {
	User_id int    `json:"id"`
	Name    string `json:"name"`
}
