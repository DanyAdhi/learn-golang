package profile

import "time"

type Profile struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Address   string     `json:"address"`
	Gender    string     `json:"gender"`
	Status    string     `json:"status"`
	CreatedAt *time.Time `json:"created_at"`
}
