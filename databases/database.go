package databases

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func DatabaseConnect() {
	db_url := os.Getenv("DATABASE_URL")

	db, err := sql.Open("mysql", db_url)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}
	fmt.Println("Database connected successfully.")
	DB = db
}
