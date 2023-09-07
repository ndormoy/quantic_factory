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

const (
	IP    = "44.333.11.22"
	login = "candidat2020"
	pwd   = "dfskj_878$*="
)

/*
To use MySQL in go program, we need to register the MySQL driver first.
*/
// func init() {
//     sql.Register("mysql", &MySQLDriver{})
// }

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

/*
Load the instance of mysql
*/
func loadMySql(login string, password string, ip string) {
	var db_key string = (login + ":" + password + "@tcp(" + ip + ")/testdb")
	db, err := sql.Open("mysql", db_key)
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}
	println("db_key  = ", db_key)
	// return db
}

func main() {
	fmt.Println("Hello, World!aaaa")
	var login string = getDotEnvVar("LOGIN")
	var password string = getDotEnvVar("PASSWORD")
	var ip string = getDotEnvVar("IP")
	fmt.Println("LOGIN: " + login)
	fmt.Println("PASSWORD: " + password)
	fmt.Println("IP: " + ip)

	loadMySql(login, password, ip)

}
