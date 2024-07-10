package initializers

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var DB *sqlx.DB

func DbConnect() {
	var Err error

	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		log.Fatal("DB_URL environment variable not set")
		return
	}
	log.Printf("Connecting to database with DSN: %s", dsn)

	DB, Err = sqlx.Connect("postgres", dsn)
	if Err != nil {
		log.Fatalf("Failed to connect to the database: %v", Err)
		return
	}

	log.Println("Successfully connected to the database")
}
