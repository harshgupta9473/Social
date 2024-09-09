package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/harshgupta9473/Social/pkg/models"
	"github.com/harshgupta9473/Social/pkg/services"
	"github.com/harshgupta9473/Social/pkg/utils"
)

type MessgaeHandler struct {
	messageService *services.MessageServices
}

func NewMsgHandler(msg *services.MessageServices) *MessgaeHandler {
	return &MessgaeHandler{messageService: msg}
}


func (handler *MessgaeHandler)HandleSendMessageToUser(w http.ResponseWriter,r *http.Request,token *jwt.Token){
	err := godotenv.Load()
	if err!=nil{
		// internal server error
		return
	}
	limit,err:=strconv.Atoi(os.Getenv("msglimit"))
	if err!=nil{
		// internal server error error to fetch data from limit
		http.Error(w,"internal server 1 ",http.StatusInternalServerError)
		return
	}
	var msgreq models.NewMessageReq
	err=json.NewDecoder(r.Body).Decode(&msgreq)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusConflict)
		return
	}
	// email,err:=utils.ExtractEmailFromJWT(token)
	id,err:=utils.ExtractIDFromJWT(token)
	if err!=nil{
		http.Error(w,"internal server 2 ",http.StatusInternalServerError)
		return
	}
	if id!=msgreq.SenderID{
		http.Error(w,"forbidden",http.StatusForbidden)
		return
	}
	err,code :=handler.messageService.SendMessage(msgreq.SenderID,msgreq.RecipentID,msgreq.Message,limit)
	if err!=nil{
		http.Error(w,err.Error(),code)
		return
	}
	utils.WriteJSON(w,http.StatusOK,"sent")
}

func (handler *MessgaeHandler) HandleGetAllRecievedMessage(w http.ResponseWriter,r *http.Request,token *jwt.Token){
	var email string
	err:=json.NewDecoder(r.Body).Decode(&email)
	if err!=nil{
		http.Error(w,"not correct formate",http.StatusBadRequest)
		return
	}
	emailFROMJWT,err:=utils.ExtractEmailFromJWT(token)
	if err!=nil{
		http.Error(w,"internal server ",http.StatusInternalServerError)
		return
	}
	if email!=emailFROMJWT{
		http.Error(w,"forbidden",http.StatusForbidden)
		return
	}
	idFromJWT,err:=utils.ExtractIDFromJWT(token)
	if err!=nil{
		http.Error(w,"internal server ",http.StatusInternalServerError)
		return
	}

	messages,err:=handler.messageService.GetAllRecievedMessages(idFromJWT)
	if err!=nil{
		http.Error(w,"error",http.StatusNoContent)
		return
	}
	utils.WriteJSON(w,http.StatusOK,messages)
	
}


func (handler *MessgaeHandler) HandleGetAllSentMessage(w http.ResponseWriter,r *http.Request,token *jwt.Token){
	var email string
	err:=json.NewDecoder(r.Body).Decode(&email)
	if err!=nil{
		http.Error(w,"not correct formate",http.StatusBadRequest)
		return
	}
	emailFROMJWT,err:=utils.ExtractEmailFromJWT(token)
	if err!=nil{
		http.Error(w,"internal server ",http.StatusInternalServerError)
		return
	}
	if email!=emailFROMJWT{
		http.Error(w,"forbidden",http.StatusForbidden)
		return
	}
	idFromJWT,err:=utils.ExtractIDFromJWT(token)
	if err!=nil{
		http.Error(w,"internal server ",http.StatusInternalServerError)
		return
	}

	messages,err:=handler.messageService.GetAllSentMessages(idFromJWT)
	if err!=nil{
		http.Error(w,"error",http.StatusNoContent)
		return
	}
	utils.WriteJSON(w,http.StatusOK,messages)
	
}
