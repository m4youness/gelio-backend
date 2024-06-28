package initializers

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var DB *sqlx.DB

func DbConnect() {
	var Err error

	dsn := os.Getenv("DB_URL")
	DB, Err = sqlx.Connect("postgres", dsn)

	if Err != nil {
		fmt.Println("Could not connect to database")
		log.Fatal(Err)

	}

}
