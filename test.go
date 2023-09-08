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
	// return fmt.Sprintf("%s:%s@tcp(%s)/%s", login, password, ip, dbName)
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?multiStatements=true", login, password, ip, dbName)
}

/*
Create the database 'quantic_db' if it doesn't exist and open the connection to it
*/
func initializeQuanticDB(login, password, ip, dbname string) (*sql.DB, error) {
	// Open the connection to 'mysql' database and create 'quantic_db'
	db, err := sql.Open("mysql", dsn(login, password, ip, "mysql"))
	if err != nil {
		return nil, err
	}

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	_, err = db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbname)
	if err != nil {
		db.Close() // Close the connection to 'mysql' database
		return nil, err
	}

	// Open the connection to 'quantic_db'
	quanticDB, err := sql.Open("mysql", dsn(login, password, ip, dbname))
	if err != nil {
		db.Close() // Close the connection to 'mysql' database
		return nil, err
	}

	return quanticDB, nil
}

/*
Create the database schema for the given Database
*/
func createDatabaseSchema(db *sql.DB) error {
	// SQL script with table creation and constraints
	schemaSQL := `

		CREATE TABLE IF NOT EXISTS ChannelType (
			ChannelTypeID smallint AUTO_INCREMENT NOT NULL,
			Name varchar(30) NOT NULL,
			PRIMARY KEY (ChannelTypeID)
		);

		CREATE TABLE IF NOT EXISTS EventType (
			EventTypeID smallint AUTO_INCREMENT NOT NULL,
			Name varchar(30) NOT NULL,
			PRIMARY KEY (EventTypeID)
		);

		CREATE TABLE IF NOT EXISTS Content (
            ContentID int AUTO_INCREMENT NOT NULL,
            ClientContentID bigint NOT NULL,
            InsertDate timestamp NOT NULL,
            PRIMARY KEY (ContentID)
        );

        CREATE TABLE IF NOT EXISTS Customer (
            CustomerID bigint NOT NULL,
            ClientCustomerID bigint NOT NULL,
            InsertDate timestamp NOT NULL,
            PRIMARY KEY (CustomerID)
        );

        CREATE TABLE IF NOT EXISTS CustomerData (
            CustomerChannelID bigint NOT NULL,
            CustomerID bigint NOT NULL,
            ChannelTypeID smallint NOT NULL,
            ChannelValue varchar(600) NOT NULL,
            InsertDate timestamp NOT NULL,
            PRIMARY KEY (CustomerChannelID),
            FOREIGN KEY (CustomerID) REFERENCES Customer (CustomerID),
            FOREIGN KEY (ChannelTypeID) REFERENCES ChannelType (ChannelTypeID)
        );

        CREATE TABLE IF NOT EXISTS CustomerEvent (
            EventID bigint NOT NULL,
            ClientEventID bigint NOT NULL,
            InsertDate timestamp NOT NULL,
            PRIMARY KEY (EventID)
        );

        CREATE TABLE IF NOT EXISTS CustomerEventData (
            EventDataId bigint NOT NULL,
            EventID bigint NOT NULL,
            ContentID int NOT NULL,
            CustomerID bigint NOT NULL,
            EventTypeID smallint NOT NULL,
            EventDate timestamp NOT NULL,
            Quantity smallint NOT NULL,
            InsertDate timestamp NOT NULL,
            PRIMARY KEY (EventDataId),
            FOREIGN KEY (EventID) REFERENCES CustomerEvent (EventID),
            FOREIGN KEY (ContentID) REFERENCES Content (ContentID),
            FOREIGN KEY (CustomerID) REFERENCES Customer (CustomerID),
            FOREIGN KEY (EventTypeID) REFERENCES EventType (EventTypeID)
        );

        CREATE TABLE IF NOT EXISTS ContentPrice (
            ContentPriceID mediumint AUTO_INCREMENT NOT NULL,
            ContentID int NOT NULL,
            Price decimal(8,2) NOT NULL,
            Currency char(3) NOT NULL,
            InsertDate timestamp NOT NULL,
            PRIMARY KEY (ContentPriceID),
            FOREIGN KEY (ContentID) REFERENCES Content (ContentID)
        );

        
    `

	_, err := db.Exec(schemaSQL)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	fmt.Println("Hello, World!aaaa")
	var login string = getDotEnvVar("LOGIN")
	var password string = getDotEnvVar("PASSWORD")
	var ip string = getDotEnvVar("IP")
	var dbname string = "quantic_db"

	quanticDB, err := initializeQuanticDB(login, password, ip, dbname)
	if err != nil {
		log.Printf("Error %s when initializing 'quantic_db' DB\n", err)
		return
	}
	defer quanticDB.Close() // Close the connection to 'quantic_db'

	// Create the database schema
	if err := createDatabaseSchema(quanticDB); err != nil {
		log.Printf("Error %s when creating database schema\n", err)
		return
	}

}
