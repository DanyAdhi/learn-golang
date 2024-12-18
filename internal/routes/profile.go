package routes

import (
	"database/sql"
	"net/http"

	"github.com/DanyAdhi/learn-golang/internal/profile"
	"github.com/DanyAdhi/learn-golang/internal/utils"
	"github.com/gorilla/mux"
)

func SetupProfileRouter(router *mux.Router, database *sql.DB) *mux.Router {
	profileRouter := router.NewRoute().Subrouter()
	profileRouter.Use(utils.AuthMiddleware)

	profileRepository := profile.NewRepository(database)
	profileService := profile.NewService(profileRepository)
	profileHandler := profile.NewHandler(profileService)

	profileRouter.HandleFunc("/profile", profileHandler.Profile).Methods(http.MethodGet)

	return profileRouter
}
