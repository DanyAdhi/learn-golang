package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/DanyAdhi/learn-golang/internal/utils"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var data RequestLogin
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Printf(`Error get body %v`, err)
		utils.ResponseError(w, http.StatusBadRequest, "Error get body")
		return
	}

	login, err := h.service.Login(data)
	if err != nil {
		if err == ErrWrongEmailOrPassword {
			utils.ResponseError(w, http.StatusBadRequest, "Email or Password wrong.")
			return
		}
		utils.ResponseError(w, http.StatusInternalServerError, "Internal server error.")
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Success", login)
}

func (h *Handler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := h.service.RefreshTokenService("refreshtoken")
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "Failed generate token")
		return
	}
	utils.ResponseSuccess(w, http.StatusOK, "success", refreshToken)
}