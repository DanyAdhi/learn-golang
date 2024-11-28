package routes

import (
	"database/sql"
	"net/http"

	"github.com/DanyAdhi/learn-golang/internal/auth"
	"github.com/gorilla/mux"
)

func SetupAuthRouter(router *mux.Router, db *sql.DB) *mux.Router {
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)

	router.HandleFunc("/sign-in", authHandler.LoginHandler).Methods(http.MethodPost)
	router.HandleFunc("/refresh-token", authHandler.RefreshTokenHandler).Methods(http.MethodPost)

	return router
}
