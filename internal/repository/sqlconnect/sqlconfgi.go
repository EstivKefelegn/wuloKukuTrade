package sqlconnect

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func ConnectDB() (*sql.DB, error) {
	fmt.Println("Connecting ... ")

	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	// user := os.Getenv("DB_USER")
	// name := os.Getenv("DB_NAME")
	// password := os.Getenv("DB_PASSWORD")
	// host := os.Getenv("DB_HOST")
	// port := os.Getenv("DB_PORT")

	// fmt.Println("USER:", user, "HOST", host, "PASSWORD", password)
	// connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, name) // "root:MyStrongPassword123!@tcp(127.0.0.1:3306)/" + dbname
	connectionString := os.Getenv("CONNECTION_STRING")
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to MariaDB")
	return db, nil
}
