package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/DanyAdhi/learn-golang/internal/utils"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h Handler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var user *UserSignUp
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Print("error body")
		utils.ResponseError(w, http.StatusBadRequest, "Error get body")
		return
	}

	message, err := utils.Validator(user)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, message)
		return
	}

	err = h.service.SignUpService(user)
	if err == ErrEmailAlreadyExist {
		log.Printf("error body %v", err)

		utils.ResponseError(w, http.StatusBadRequest, "Email already exist")
		return
	}
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	utils.ResponseSuccess(w, http.StatusCreated, "Success", nil)
}

func (h *Handler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	var data RequestSignIn
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Printf(`Error get body %v`, err)
		utils.ResponseError(w, http.StatusBadRequest, "Error get body")
		return
	}

	signIn, err := h.service.SignIn(data)
	if err != nil {
		if err == ErrWrongEmailOrPassword {
			utils.ResponseError(w, http.StatusBadRequest, "Email or Password wrong.")
			return
		}
		utils.ResponseError(w, http.StatusInternalServerError, "Internal server error.")
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Success", signIn)
}

func (h *Handler) SignOutHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(utils.UserKey).(*utils.JwtDecodeInterface)
	token := getBearerToken(r)

	err := h.service.SignOutService(user.ID, token)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "internal server error.")
		return
	}
	utils.ResponseSuccess(w, http.StatusOK, "Success", nil)
}

func (h *Handler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var bodyRefreshToken ReqBodyRefreshToken

	err := json.NewDecoder(r.Body).Decode(&bodyRefreshToken)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "Failed get data body")
		return
	}

	refreshToken, err := h.service.RefreshTokenService(bodyRefreshToken.Refresh_token)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "Failed generate token")
		return
	}
	utils.ResponseSuccess(w, http.StatusOK, "success", refreshToken)
}

func getBearerToken(r *http.Request) string {
	authorization := r.Header.Get("Authorization")
	if authorization == "" || !strings.Contains(authorization, "Bearer") {
		return ""
	}
	token := strings.Split(authorization, " ")[1]
	return token
}
