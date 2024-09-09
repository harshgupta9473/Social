package utils

import (
	"database/sql"

	"github.com/harshgupta9473/Social/configs"
	_ "github.com/lib/pq"
)

func GetDB() *sql.DB {
	return configs.DB
}
