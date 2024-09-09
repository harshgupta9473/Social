package createtable

import (
	"database/sql"
	_ "github.com/lib/pq"
)




func CreateUserTable(DB *sql.DB) error {
	query := `create table if not exists users(
	id integer generated always as identity primary key,
	email varchar(255) not null,
	encrypted_password varchar(100) not null,
	created_at timestamp not null,
	verified boolean default true
	)`
	_, err := DB.Exec(query)
	return err
}

func CreateTempUserTable(DB *sql.DB) error {
	query := `create table if not exists tempusers(
	id integer generated always as identity primary key,
	email varchar(255) not null,
	encrypted_password varchar(100) not null,
	token varchar(32) not null,
	expires_at timestamp not null
	)`
	_, err := DB.Exec(query)
	return err
}

func CreateUserProfileTable(db *sql.DB) error {

	query := `create table if not exists userprofiles(
	id integer generated always as identity primary key,
	username varchar(100) unique not null,
	email varchar(255) unique not null,
	firstname varchar(100) not null,
	lastname varchar(100),
	bio varchar(255),
	interests jsonb,
	location varchar(255),
	created_at timestamp not null default current_timestamp ,
	updated_at timestamp not null default current_timestamp
	)`
	_, err := db.Exec(query)
	return err
}

func CreateMessageTable(db *sql.DB)error{
	query:=`create table if not exists messagetable(
	id integer generated always as identity primary key,
	sender_id integer not null,
	recipent_id integer not null,
	message text not null,
	sent_at timestamp not null default current_timestamp,
	foreign key (sender_id) references userprofiles(id),
	foreign key (recipent_id) references userprofiles(id)
	)`
	_,err:=db.Exec(query)
	return err
}

func CreateUserMessageCountTable(db *sql.DB)error{
	query:=`create table if not exists messagecounttable(
	user_id integer not null,
	message_date  date not null,
	message_count integer not null default 0,
	primary key (user_id,message_date),
	foreign key (user_id) references userprofiles(id)
	)`
	_,err:=db.Exec(query)
	return err
}



func TableInit(db *sql.DB) error{
	err:=CreateUserTable(db)
	if err!=nil{
		return err
	}
	err=CreateUserProfileTable(db)
	if err!=nil{
		return err
	}
	err=CreateTempUserTable(db)
	if err!=nil{
		return err
	}
	err=CreateMessageTable(db)
	if err!=nil{
		return err
	}
	err=CreateUserMessageCountTable(db)
	if err!=nil{
		return err
	}
	return nil
}