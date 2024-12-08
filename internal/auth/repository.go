package auth

import (
	"database/sql"
	"log"
)

type Repository interface {
	StoreUsersSignUpRepository(user *UserSignUp) error
	GetUsersByEmail(email string) (*User, error)
	StoreRefreshToken(userId int, refreshToken string) error
	GetRefreshToken(refreshToken string) (*GetRefreshToken, error)
	RevokeToken(userId int) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) StoreUsersSignUpRepository(user *UserSignUp) error {
	_, err := r.db.Exec(
		`INSERT INTO users (name, email, password) VALUES ($1, $2, $3)`,
		user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) GetUsersByEmail(email string) (*User, error) {
	row := r.db.QueryRow(`SELECT id, name, password FROM users WHERE email = $1 AND deletedat is NULL`, email)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repository) StoreRefreshToken(userId int, refreshToken string) error {
	var id int
	err := r.db.QueryRow(`SELECT id FROM refresh_tokens WHERE user_id = $1`, userId).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if id == 0 {
		r.db.Exec(
			`INSERT INTO refresh_tokens(user_id, token, is_revoked) VALUES ($1, $2, $3)`,
			userId, refreshToken, false,
		)
	} else {
		r.db.Exec(
			`UPDATE refresh_tokens SET token = $1, is_revoked = $2 WHERE id = $3`,
			refreshToken, false, id,
		)
	}
	return nil
}

func (r *repository) GetRefreshToken(refreshToken string) (*GetRefreshToken, error) {
	row := r.db.QueryRow(
		`SELECT refresh_tokens.user_id, users.name FROM refresh_tokens 
    	LEFT JOIN users ON users.id = refresh_tokens.user_id
      WHERE refresh_tokens.is_revoked = $1 AND refresh_tokens.token = $2`,
		false, refreshToken,
	)
	var dataRefreshToken GetRefreshToken
	err := row.Scan(&dataRefreshToken.User_id, &dataRefreshToken.Name)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		log.Printf("error get refresh token")
		return nil, err
	}
	return &dataRefreshToken, nil
}

func (r *repository) RevokeToken(userId int) error {
	_, err := r.db.Exec(`UPDATE refresh_tokens SET is_revoked = $1 WHERE user_id = $2`, true, userId)
	if err != nil {
		return err
	}
	return nil
}
