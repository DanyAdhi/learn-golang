package routes

import (
	"database/sql"
	"net/http"

	"github.com/DanyAdhi/learn-golang/internal/users"
	"github.com/DanyAdhi/learn-golang/internal/utils"
	"github.com/gorilla/mux"
)

func SetupUserRouter(router *mux.Router, database *sql.DB) *mux.Router {
	usersRepo := users.NewRepository(database)
	usersService := users.NewService(usersRepo)
	userHandler := users.NewHandler(usersService)

	router.Handle("/users", utils.AuthMiddleware(http.HandlerFunc(userHandler.GetAllUsersHandler))).Methods(http.MethodGet)
	router.Handle("/users/{id}", utils.AuthMiddleware(http.HandlerFunc(userHandler.GetOneUsersHandler))).Methods(http.MethodGet)
	router.Handle("/users", utils.AuthMiddleware(http.HandlerFunc(userHandler.CreateUsersHandler))).Methods(http.MethodPost)
	router.Handle("/users/{id}", utils.AuthMiddleware(http.HandlerFunc(userHandler.UpdateUsersHandler))).Methods(http.MethodPut)
	router.Handle("/users/{id}", utils.AuthMiddleware(http.HandlerFunc(userHandler.DeleteusersHandler))).Methods(http.MethodDelete)

	return router
}
