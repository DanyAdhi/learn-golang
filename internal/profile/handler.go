package profile

import (
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

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {

	user, ok := r.Context().Value(utils.UserKey).(*utils.JwtDecodeInterface)
	if !ok {
		log.Print("Error get data id in context")
		utils.ResponseError(w, http.StatusInternalServerError, "Failed get profil")
		return
	}

	profile, err := h.service.GetProfileService(user.ID)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "Failed get profil")
		return
	}
	utils.ResponseSuccess(w, http.StatusOK, "Success", profile)
}
