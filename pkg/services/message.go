package services

import (
	"fmt"
	"log"
	"net/http"

	"github.com/harshgupta9473/Social/pkg/models"
	"github.com/harshgupta9473/Social/pkg/repository"
)

type MessageServices struct {
	MessageRepo *repository.MessageRepository
}

func NewMessageService(userRepo *repository.MessageRepository)* MessageServices{
	return &MessageServices{MessageRepo: userRepo}
}


func (services *MessageServices) SendMessage(senderID ,recipetID uint64,content string,limit int)(error,int){
	count,err:=services.MessageRepo.GetDailyMsgCount(senderID)
	if err!=nil{
		return err,http.StatusInternalServerError
	}
	if count==limit{
		return fmt.Errorf("daily limit reached"),http.StatusConflict
	}else{
		err=services.MessageRepo.SendMessage(senderID,recipetID,content)
		if err!=nil{
			log.Println(err)
			return fmt.Errorf("not able to send message"),http.StatusInternalServerError
		}
		err=services.MessageRepo.IncrementMsgCount(senderID)
		if err!=nil{
			log.Println(err)
			return fmt.Errorf("not able to increment msgcount"),http.StatusInternalServerError
		}
	}
	return nil,http.StatusOK
	
}

func (services *MessageServices) GetAllRecievedMessages(userID uint64)([]models.Message,error){
	messages,err:=services.MessageRepo.GetAllRecievedMessages(userID)
	if err!=nil{
		return nil,err
	}
	return messages,nil
}

func (services *MessageServices) GetAllSentMessages(userID uint64)([]models.Message,error){
	messages,err:=services.MessageRepo.GetAllSentMessages(userID)
	if err!=nil{
		return nil,err
	}
	return messages,nil
}


func (services *MessageServices) GetAllMessages(userID uint64)([]models.Message,error){
	messages,err:=services.MessageRepo.GetAllMessages(userID)
	if err!=nil{
		return nil,err
	}
	return messages,nil
}

func (services *MessageServices)GetDailyMsgCount(userID uint64)(int ,error){
	count,err:=services.MessageRepo.GetDailyMsgCount(userID)
	if err!=nil{
		return 0,err
	}
	return count,nil
}
