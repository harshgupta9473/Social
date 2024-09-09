package repository

import (
	"database/sql"
	"time"

	"github.com/harshgupta9473/Social/pkg/models"
)
type MessageRepository struct{
	DB *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{DB: db}
}

func (repo *MessageRepository) SendMessage(senderID uint64, recipentID uint64, content string) error {
	query := `insert into messagetable
	(sender_id,recipent_id,message,sent_at)
	values($1,$2,$3,$4)`
	_, err := repo.DB.Exec(query, senderID, recipentID, content, time.Now().UTC())
	if err!=nil{
		return err
	}
	return nil
}

func (repo *MessageRepository) IncrementMsgCount(senderID uint64)error{
	query:=`insert into messagecounttable
	(user_id,message_date,message_count)
	values($1,$2,1)
	on conflict (user_id,message_date)
	do update set message_count=messagecounttable.message_count+1`
	_,err:=repo.DB.Exec(query,senderID,time.Now().Format("2006-01-02"))
	if err!=nil{
		return err
	}
	return nil
}

func (repo *MessageRepository) GetDailyMsgCount(senderID uint64)(int ,error){
	query:=`select message_count from messagecounttable where user_id=$1 and message_date=$2`
	rows:=repo.DB.QueryRow(query,senderID,time.Now().Format("2006-01-02"))
	var count int
		err:=rows.Scan(&count)
		if err!=nil{
			if err==sql.ErrNoRows{
				return 0,nil
			}
			return 0,err
		}
	return count,nil
}

func (repo *MessageRepository) GetAllSentMessages(userID uint64)([]models.Message,error){
	query:=`select * from messagetable where sender_id=$1`
	rows,err:=repo.DB.Query(query,userID)
	if err!=nil{
		return nil,err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next(){
		var message models.Message
		err=rows.Scan(&message.ID,&message.SenderID,&message.RecipientID,&message.Content,&message.SentAt);
		if err!=nil{
			return nil,err
		}
		messages=append(messages, message)
	}
	return messages,nil
}

func (repo *MessageRepository)GetAllRecievedMessages(userID uint64)([]models.Message,error){
	query:=`select * from messagetable where recipent_id=$1`
	rows,err:=repo.DB.Query(query,userID)
	if err!=nil{
		return nil,err
	}
	defer rows.Close()

	var messages []models.Message

	for rows.Next(){
		var message models.Message
		err=rows.Scan(&message.ID,&message.SenderID,&message.RecipientID,&message.Content,&message.SentAt)
		if err!=nil{
			return nil,err
		}
		messages=append(messages,message)
	}
	return messages,nil
}

func (repo *MessageRepository)GetAllMessages (userID uint64)([]models.Message,error){
	query:=`select * from messagetable where sender_id=$1 or recipent_id=$1`
	rows,err:=repo.DB.Query(query,userID)
	if err!=nil{
		return nil,err
	}
	var messages []models.Message
	for rows.Next(){
		var message models.Message
		err=rows.Scan(&message.ID,&message.SenderID,&message.RecipientID,&message.SentAt)
		if err!=nil{
			return nil,err
		}
		messages=append(messages, message)
	}
	return messages,nil
}