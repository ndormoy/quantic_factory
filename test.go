package main

/*
Import github.com/go-sql-driver/mysql
Permit to import MySQL
*/
import (
	"fmt"
	"log"
	"os"

	// Import godotenv
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Customer struct {
}

/*
Function to get back the var in .env file
*/
func getDotEnvVar(key string) string {

	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func main() {
	fmt.Println("Hello, World!aaaa")
	var login string = getDotEnvVar("LOGIN")
	var password string = getDotEnvVar("PASSWORD")
	var ip string = getDotEnvVar("IP")
	fmt.Println("LOGIN: " + login)
	fmt.Println("PASSWORD: " + password)
	fmt.Println("IP: " + ip)

	//Creating the MySQL instance
	var db_key string = (login + ":" + password + "@tcp(" + ip + ")/testdb")
	db, err := sql.Open("mysql", db_key)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	fmt.Println("Success!, Database MySQL is connected")
}
