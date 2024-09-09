package main

import (
	"log"
	"net/http"

	"github.com/harshgupta9473/Social/configs"
	createtable "github.com/harshgupta9473/Social/pkg/createTable"
	"github.com/harshgupta9473/Social/pkg/handlers"
	"github.com/harshgupta9473/Social/pkg/repository"
	"github.com/harshgupta9473/Social/pkg/routes"
	"github.com/harshgupta9473/Social/pkg/services"
	"github.com/harshgupta9473/Social/pkg/utils"
)

func main() {
	configs.InitDB()

	db:=utils.GetDB()
	err:=createtable.TableInit(db)
	if err!=nil{
		log.Fatal(err)
	}

	userRepo:=repository.NewUserRepository(db)
	authService:=services.NewAuthService(userRepo)
	userHandler:=handlers.NewUserHandler(authService)
	msgRepo:=repository.NewMessageRepository(db)
	msgService:=services.NewMessageService(msgRepo)
	msgHandler:=handlers.NewMsgHandler(msgService)


	router:=routes.SetupRoutes(userHandler,msgHandler)
	http.ListenAndServe(":3000",router)
	
}