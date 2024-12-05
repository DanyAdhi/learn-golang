package users

import (
	"database/sql"
	"log"
)

type Repository interface {
	GetAllUsersRepository(params GetAllUsersParmas) (*[]User, int, error)
	GetOneUsersRepository(id int) (*User, error)
	StoreUsersRepository(user *Createuser, password string) error
	CheckEmailExists(email string) (bool, error)
	UpdateUsersRepository(id int, user *UpdateUser) error
	DeleteUsersRepository(id int) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAllUsersRepository(p GetAllUsersParmas) (*[]User, int, error) {

	var totalRecords int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&totalRecords)
	if err != nil {
		log.Printf("Error get total records user. %v", err)
		return nil, 0, err
	}

	offset := (p.Page - 1) * p.Limit

	rows, err := r.db.Query(
		`SELECT id, name, email, address, gender FROM users ORDER BY id DESC LIMIT $1 OFFSET $2`,
		p.Limit, offset,
	)
	if err != nil {
		log.Printf("Error get data user: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var dataUser User
		err := rows.Scan(&dataUser.ID, &dataUser.Name, &dataUser.Email, &dataUser.Address, &dataUser.Gender)
		if err != nil {
			log.Printf("Row scan error: %v", err)
			return nil, 0, err
		}
		users = append(users, dataUser)
	}

	return &users, totalRecords, nil
}

func (r *repository) GetOneUsersRepository(id int) (*User, error) {
	row := r.db.QueryRow(`SELECT id, name, email, address, gender FROM users WHERE id = $1`, id)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Address, &user.Gender)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) CheckEmailExists(email string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(1) FROM users WHERE email = $1`, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *repository) StoreUsersRepository(user *Createuser, password string) error {
	_, err := r.db.Exec(
		`INSERT INTO users (name, email, address, gender, password) VALUES ($1, $2, $3, $4, $5)`,
		user.Name, user.Email, user.Address, user.Gender, password)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) UpdateUsersRepository(id int, user *UpdateUser) error {
	_, err := r.db.Exec(
		`UPDATE users SET name = $1, address = $2, gender = $3 WHERE id = $4`,
		user.Name, user.Address, user.Gender, id,
	)

	if err != nil {
		return err
	}
	return nil
}

func (r *repository) DeleteUsersRepository(id int) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}
