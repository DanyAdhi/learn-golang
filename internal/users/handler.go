package users

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/DanyAdhi/learn-golang/internal/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAllUsersService()
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "Failed get data users")
		return
	}
	utils.ResponseSuccess(w, http.StatusOK, "OK", users)
}

func (h *Handler) GetOneUsersHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Id not valid.")
		return
	}

	user, err := h.service.GetOneUsersService(id)
	if err != nil {
		if err == ErrUserNotFound {
			utils.ResponseError(w, http.StatusNotFound, "Data user not found")
			return
		}
		utils.ResponseError(w, http.StatusInternalServerError, "Failed get data users")
		return
	}
	utils.ResponseSuccess(w, http.StatusOK, "OK", user)
}

func (h *Handler) CreateUsersHandler(w http.ResponseWriter, r *http.Request) {
	var user Createuser
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Error get body %v", err)
		utils.ResponseError(w, http.StatusBadRequest, "Error get data body")
		return
	}

	err = h.service.CreateUsersService(&user)
	if err == ErrEmailAlreadyExists {
		utils.ResponseError(w, http.StatusBadRequest, "Email already exists")
		return
	}
	if err == ErrGeneratePassword {
		utils.ResponseError(w, http.StatusBadRequest, "Failed generate password")
		return
	}
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "Failed create user.")
		return
	}

	utils.ResponseSuccess(w, http.StatusCreated, "Created", nil)
}

func (h *Handler) UpdateUsersHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Id not valid.")
		return
	}

	var user UpdateUser
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Error get data body. %v", err)
		utils.ResponseError(w, http.StatusBadRequest, "Failed get data body")
		return
	}

	err = h.service.UpdateUsersService(id, &user)
	if err == ErrUserNotFound {
		utils.ResponseError(w, http.StatusNotFound, "User not found.")
		return
	}
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "Failed update user")
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Success", nil)
}

func (h *Handler) DeleteusersHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("Vars %v", vars)
	utils.ResponseSuccess(w, http.StatusOK, "Success", nil)
}
