package configs

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	connstr := os.Getenv("connStr")
	db, err := sql.Open("postgres", connstr)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	DB = db
	log.Println("Database Conneection established")
}


