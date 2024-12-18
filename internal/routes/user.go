package routes

import (
	"database/sql"
	"net/http"

	"github.com/DanyAdhi/learn-golang/internal/users"
	"github.com/DanyAdhi/learn-golang/internal/utils"
	"github.com/gorilla/mux"
)

func SetupUserRouter(router *mux.Router, database *sql.DB) *mux.Router {
	userRouter := router.NewRoute().Subrouter()
	userRouter.Use(utils.AuthMiddleware)

	usersRepo := users.NewRepository(database)
	usersService := users.NewService(usersRepo)
	userHandler := users.NewHandler(usersService)

	userRouter.HandleFunc("/users", userHandler.GetAllUsersHandler).Methods(http.MethodGet)
	userRouter.HandleFunc("/users/{id}", userHandler.GetOneUsersHandler).Methods(http.MethodGet)
	userRouter.HandleFunc("/users", userHandler.CreateUsersHandler).Methods(http.MethodPost)
	userRouter.HandleFunc("/users/{id}", userHandler.UpdateUsersHandler).Methods(http.MethodPut)
	userRouter.HandleFunc("/users/{id}", userHandler.DeleteusersHandler).Methods(http.MethodDelete)

	return router
}
