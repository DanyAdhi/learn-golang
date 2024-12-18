package profile

import "database/sql"

type Repository interface {
	GetProfileRepository(id int) (*Profile, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetProfileRepository(id int) (*Profile, error) {
	row := r.db.QueryRow(
		`SELECT id, name, email, address, gender, status, createdat FROM users WHERE id = $1`,
		id,
	)

	var profile Profile
	err := row.Scan(&profile.ID, &profile.Name, &profile.Email, &profile.Address, &profile.Gender, &profile.Status, &profile.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}
