package main

/*
Import github.com/go-sql-driver/mysql
Permit to import MySQL
*/
import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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

func dsn(login string, password string, ip string, dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", login, password, ip, dbName)
}

func main() {
	fmt.Println("Hello, World!aaaa")
	var login string = getDotEnvVar("LOGIN")
	var password string = getDotEnvVar("PASSWORD")
	var ip string = getDotEnvVar("IP")
	var dbname string = "quantic_db"
	fmt.Println("LOGIN: " + login)
	fmt.Println("PASSWORD: " + password)
	fmt.Println("IP: " + ip)

	//Creating the MySQL instance
	// var db_key string = (login + ":" + password + "@tcp(" + ip + ")/testdb")
	// db, err := sql.Open("mysql", db_key)
	// if err != nil {
	// 	panic(err.Error())
	// }
	// defer db.Close()
	// fmt.Println("Success!, Database MySQL is connected")
	db, err := sql.Open("mysql", dsn(login, password, ip, "mysql"))
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return
	}

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	// res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbname)
	res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbname)
	if err != nil {
		log.Printf("Error %s when creating DB\n", err)
		fmt.Println(res)
		return
	}
	db.Close() // Close the connection to 'mysql' database

	db, err = sql.Open("mysql", dsn(login, password, ip, dbname))
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return
	}
	defer db.Close()

}
