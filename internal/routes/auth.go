package routes

import (
	"database/sql"
	"net/http"

	"github.com/DanyAdhi/learn-golang/internal/auth"
	"github.com/DanyAdhi/learn-golang/internal/utils"
	"github.com/gorilla/mux"
)

func SetupAuthRouter(router *mux.Router, db *sql.DB) *mux.Router {
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)

	router.HandleFunc("/sign-in", authHandler.SignInHandler).Methods(http.MethodPost)
	router.HandleFunc("/refresh-token", authHandler.RefreshTokenHandler).Methods(http.MethodPost)
	router.Handle("/sign-out", utils.AuthMiddleware(http.HandlerFunc(authHandler.SignOutHandler))).Methods(http.MethodPost)

	return router
}
