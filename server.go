package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
	Gender  string `json:"gender"`
}

type Meta struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResponseSuccess struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

var db *sql.DB

func main() {
	var err error

	conectionString := "postgres://root:root1234@localhost:5439/learn?sslmode=disable"
	db, err = sql.Open("postgres", conectionString)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Cannot connect to databbase: %v", err)
	}

	fmt.Println("Success connect database.")

	http.HandleFunc("/users", handleUsers)
	http.HandleFunc("/users/", handleSingleUsers)

	fmt.Println("Server is running on port 3001")
	log.Fatal(http.ListenAndServe(":3001", nil))
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getAllUsers(w)
	case http.MethodPost:
		storeUsers(w, r)
	default:
		responseError(w, http.StatusMethodNotAllowed, "Method now allowed")
	}
}

func handleSingleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getOneUsers(w, r)
	case http.MethodPut:
		updateUsers(w, r)
	case http.MethodDelete:
		deleteUsers(w, r)
	default:
		responseError(w, http.StatusMethodNotAllowed, "Method now allowed")
	}
}

func getAllUsers(w http.ResponseWriter) {
	rows, err := db.Query("SELECT id, name, email, address, gender FROM users ORDER BY id DESC LIMIT 10")
	if err != nil {
		log.Printf("Error get data user: %v", err)
		responseError(w, http.StatusInternalServerError, "Failed get data users")
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var dataUser User
		err := rows.Scan(&dataUser.ID, &dataUser.Name, &dataUser.Email, &dataUser.Address, &dataUser.Gender)
		if err != nil {
			log.Printf("Row scan error: %v", err)
			responseError(w, http.StatusInternalServerError, "")
		}

		users = append(users, dataUser)
	}

	responseSuccess(w, http.StatusOK, "OK", users)
}

func getOneUsers(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")

	id, err := strconv.Atoi(path[len(path)-1])
	if err != nil {
		responseError(w, http.StatusBadRequest, "Id Not falid.")
		return
	}

	row := db.QueryRow("SELECT id, name, email, address, gender FROM users where id =$1", id)

	var user User
	err = row.Scan(&user.ID, &user.Name, &user.Email, &user.Address, &user.Gender)
	if err == sql.ErrNoRows {
		responseError(w, http.StatusNotFound, "Data user not found.")
		return
	}
	if err != nil {
		log.Printf("Failed scan user: %v", err)
		responseError(w, http.StatusInternalServerError, "Failed get data user")
		return
	}

	responseSuccess(w, http.StatusOK, "Success", user)
}

func storeUsers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		log.Printf("Error get body. %v", err)
		responseError(w, http.StatusBadRequest, "Error get body")
		return
	}

	var checkEmail string
	err = db.QueryRow(`SELECT email FROM users WHERE email = $1`, newUser.Email).Scan(&checkEmail)
	if err == nil {
		responseError(w, http.StatusBadRequest, "Email already exists.")
		return
	} else if err != sql.ErrNoRows {
		// Handle unexpected errors
		log.Printf("Error checking email: %v", err)
		responseError(w, http.StatusInternalServerError, "Failed to check email.")
		return
	}

	hashPassword, err := hashPassword("password")
	if err != nil {
		responseError(w, http.StatusBadRequest, "Failed generate default password")
	}

	_, err = db.Exec(
		"INSERT INTO users (name, email, address, gender, password) VALUES ($1, $2, $3, $4, $5)",
		newUser.Name,
		newUser.Email,
		newUser.Address,
		newUser.Gender,
		hashPassword,
	)
	if err != nil {
		log.Printf("Error insert user. %v", err)
		responseError(w, http.StatusInternalServerError, "Failed insert user.")
		return
	}

	responseSuccess(w, http.StatusCreated, "Success", nil)
}

func updateUsers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		log.Printf("Error get body. %v", err)
		responseError(w, http.StatusBadRequest, "Error get body")
		return
	}

	id, err := getId(r)
	if err != nil {
		responseError(w, http.StatusBadRequest, "Id not valid.")
		return
	}

	var user User
	err = db.QueryRow(`SELECT id FROM users WHERE id = $1`, id).Scan(&user.ID)
	if err != nil {
		responseError(w, http.StatusBadRequest, "Data user not found.")
		return
	}

	_, err = db.Exec(
		`UPDATE users SET name = $1, gender = $2, address = $3 WHERE id = $4`,
		newUser.Name,
		newUser.Gender,
		newUser.Address,
		id,
	)
	if err != nil {
		log.Printf("Error update user %d. %v", id, err)
		responseError(w, http.StatusInternalServerError, "Failed update users")
		return
	}

	responseSuccess(w, http.StatusOK, "Success", nil)
}

func deleteUsers(w http.ResponseWriter, r *http.Request) {
	id, err := getId(r)
	if err != nil {
		responseError(w, http.StatusBadRequest, "Id not valid.")
		return
	}

	_, err = db.Exec(
		`DELETE FROM users WHERE id = $1`,
		id,
	)
	if err != nil {
		log.Printf("Error delete user %d. %v", id, err)
		responseError(w, http.StatusInternalServerError, "Failed delete user.")
		return
	}

	responseSuccess(w, http.StatusOK, "Success", nil)
}

func responseSuccess(w http.ResponseWriter, code int, message string, data interface{}) {
	if message == "" {
		message = "success"
	}

	if code < 100 || code >= 600 {
		code = 200
	}

	response := ResponseSuccess{
		Meta: Meta{
			Success: true,
			Code:    code,
			Message: message,
		},
		Data: data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

func responseError(w http.ResponseWriter, code int, message string) {
	if message == "" {
		message = "Internal server error"
	}

	if code < 100 || code >= 600 {
		code = 500
	}

	response := Meta{
		Success: false,
		Code:    code,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

func hashPassword(pasword string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(pasword), 10)
	if err != nil {
		log.Printf("Error hash password. %v", err)
		return "", err
	}
	return string(hashPassword), nil
}

func getId(r *http.Request) (int, error) {
	path := r.URL.Path
	split := strings.Split(path, `/`)
	id, err := strconv.Atoi(split[len(split)-1])
	if err != nil {
		log.Printf("Error get id. %v", err)
		return 0, err
	}
	return id, nil
}
