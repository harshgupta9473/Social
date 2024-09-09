package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/harshgupta9473/Social/pkg/handlers"
	"github.com/harshgupta9473/Social/pkg/middleware"
)

func SetupRoutes(userHandler *handlers.UserHandler, msgHandler *handlers.MessgaeHandler)(*mux.Router){
	router:=mux.NewRouter()

	router.HandleFunc("/register",userHandler.HandleRegistration).Methods(http.MethodPost)
	router.HandleFunc("/verify",userHandler.HandleVerification).Methods(http.MethodGet)

	router.HandleFunc("/login",userHandler.HandleUserLogin).Methods(http.MethodPost)

	router.HandleFunc("/profile",middleware.WithJWTAuth(userHandler.GetUserProfile))
	router.HandleFunc("/profile/create",middleware.WithJWTAuth(userHandler.HandleCreateUserProfile)).Methods(http.MethodPost)
	router.HandleFunc("/profile/update",middleware.WithJWTAuth(userHandler.HandleUserProfileUpdation)).Methods(http.MethodPost)

	router.HandleFunc("/message",middleware.WithJWTAuth(msgHandler.HandleSendMessageToUser)).Methods(http.MethodPost)
	router.HandleFunc("/messages/recieved",middleware.WithJWTAuth(msgHandler.HandleGetAllRecievedMessage)).Methods(http.MethodGet,http.MethodPost)
	router.HandleFunc("/messages/sent",middleware.WithJWTAuth(msgHandler.HandleGetAllSentMessage)).Methods(http.MethodGet,http.MethodPost)

	return router
}