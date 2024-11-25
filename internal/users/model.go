package users

type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
	Gender  string `json:"gender"`
}

type Createuser struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
	Gender  string `json:"gender"`
}

type UpdateUser struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Gender  string `json:"gender"`
}
