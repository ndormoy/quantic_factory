package main

/*
Import github.com/go-sql-driver/mysql
Permit to import MySQL
*/
import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"strings"
	"time"

	// "bufio"

	// Import godotenv
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/schollz/progressbar/v3"
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
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?multiStatements=true&clientFoundRows=true&allowAllFiles=true", login, password, ip, dbName)
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
			ChannelTypeID smallint UNSIGNED AUTO_INCREMENT NOT NULL,
			Name varchar(30) NOT NULL,
			PRIMARY KEY (ChannelTypeID)
		);

		CREATE TABLE IF NOT EXISTS EventType (
			EventTypeID smallint UNSIGNED AUTO_INCREMENT NOT NULL,
			Name varchar(30) NOT NULL,
			PRIMARY KEY (EventTypeID)
		);

		CREATE TABLE IF NOT EXISTS Content (
            ContentID int UNSIGNED AUTO_INCREMENT NOT NULL,
            ClientContentID bigint UNSIGNED NOT NULL,
            InsertDate timestamp NOT NULL,
            PRIMARY KEY (ContentID)
        );

        CREATE TABLE IF NOT EXISTS Customer (
            CustomerID bigint UNSIGNED AUTO_INCREMENT NOT NULL,
            ClientCustomerID bigint UNSIGNED NOT NULL,
            InsertDate timestamp NOT NULL,
            PRIMARY KEY (CustomerID)
        );

        CREATE TABLE IF NOT EXISTS CustomerData (
            CustomerChannelID bigint UNSIGNED AUTO_INCREMENT NOT NULL,
            CustomerID bigint UNSIGNED NOT NULL,
            ChannelTypeID smallint UNSIGNED NOT NULL,
            ChannelValue varchar(600) NOT NULL,
            InsertDate timestamp NOT NULL,
            PRIMARY KEY (CustomerChannelID),
            FOREIGN KEY (CustomerID) REFERENCES Customer (CustomerID),
            FOREIGN KEY (ChannelTypeID) REFERENCES ChannelType (ChannelTypeID)
        );

        CREATE TABLE IF NOT EXISTS CustomerEvent (
            EventID bigint UNSIGNED AUTO_INCREMENT NOT NULL,
            ClientEventID bigint NOT NULL,
            InsertDate timestamp NOT NULL,
            PRIMARY KEY (EventID)
        );

        CREATE TABLE IF NOT EXISTS CustomerEventData (
            EventDataId bigint UNSIGNED AUTO_INCREMENT NOT NULL,
            EventID bigint UNSIGNED NOT NULL,
            ContentID int UNSIGNED NOT NULL,
            CustomerID bigint UNSIGNED NOT NULL,
            EventTypeID smallint UNSIGNED NOT NULL,
            EventDate timestamp NOT NULL,
            Quantity smallint UNSIGNED NOT NULL,
            InsertDate timestamp NOT NULL,
            PRIMARY KEY (EventDataId),
            FOREIGN KEY (EventID) REFERENCES CustomerEvent (EventID),
            FOREIGN KEY (ContentID) REFERENCES Content (ContentID),
            FOREIGN KEY (CustomerID) REFERENCES Customer (CustomerID),
            FOREIGN KEY (EventTypeID) REFERENCES EventType (EventTypeID)
        );

        CREATE TABLE IF NOT EXISTS ContentPrice (
            ContentPriceID mediumint UNSIGNED AUTO_INCREMENT NOT NULL,
            ContentID int UNSIGNED NOT NULL,
            Price decimal(8,2) UNSIGNED NOT NULL,
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

/*
This function insert "rows" in the given table "tableName"
*/

func insertRowsInDb(db *sql.DB, rows []string, tableName string) error {
	for _, row := range rows {
		query := fmt.Sprintf("INSERT INTO %s (Name) VALUES (?)", tableName)
		_, err := db.Exec(query, row)
		if err != nil {
			return err
		}
	}

	return nil
}

/*
This function creates the tables ChannelType and EventType
*/
func createTypesTables(db *sql.DB) error {
	channelTypes := []string{"Email", "PhoneNumber", "Postal", "MobileID", "Cookie"}
	if err := insertRowsInDb(db, channelTypes, "ChannelType"); err != nil {
		log.Printf("Error %s when inserting ChannelTypes\n", err)
		return err
	}
	eventTypes := []string{"sent", "view", "click", "visit", "cart", "purchase"}
	if err := insertRowsInDb(db, eventTypes, "EventType"); err != nil {
		log.Printf("Error %s when inserting EventTypes\n", err)
		return err
	}
	return nil
}

// Function to import data from a CSV file and insert it into a specified table, with any number of rows "columnNames"
func importDataFromCSV(db *sql.DB, csvPath, tableName string, columnNames []string) error {
	// Open the CSV file
	file, err := os.Open(csvPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read CSV records
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Prepare the SQL statement for insertion
	numColumns := len(columnNames)
	placeholders := make([]string, numColumns)
	for i := 0; i < numColumns; i++ {
		placeholders[i] = "?"
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, strings.Join(columnNames, ", "), strings.Join(placeholders, ", "))
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Skip the first row (header)
	for i, record := range records {
		if i == 0 { // Skip the first row
			continue
		}

		// Check if the number of columns in the record matches the expected number
		if len(record) != numColumns {
			return fmt.Errorf("record at row %d has an incorrect number of columns", i+1)
		}

		// Convert the record values to interface{} for Exec
		var values []interface{}
		for _, colValue := range record {
			values = append(values, colValue)
		}

		_, err := stmt.Exec(values...)
		if err != nil {
			return err
		}
	}

	return nil
}

/*
This function fills all the tables with the data from the CSV files
*/

func fillAllTables(db *sql.DB, progressBar *progressbar.ProgressBar) error {
	//Creates columnNames for each table in the schema
	columnNamesCustomer := []string{"CustomerID", "ClientCustomerID", "InsertDate"}
	columnNamesCustomerData := []string{"CustomerChannelID", "CustomerID", "ChannelTypeID", "ChannelValue", "InsertDate"}
	columnNamesCustomerEvent := []string{"EventID", "ClientEventID", "InsertDate"}
	columnNamesContent := []string{"ContentID", "ClientContentID", "InsertDate"}
	columnNamesContentPrice := []string{"ContentPriceID", "ContentID", "Price", "Currency", "InsertDate"}
	columNamesCustomerEventData := []string{"EventDataID", "EventID", "ContentID", "CustomerID", "EventTypeID", "EventDate", "Quantity", "InsertDate"}

	steps := []struct {
		csvPath     string
		tableName   string
		columnNames []string
	}{
		{"csv/Customer.csv", "Customer", columnNamesCustomer},
		{"csv/CustomerData.csv", "CustomerData", columnNamesCustomerData},
		{"csv/CustomerEvent.csv", "CustomerEvent", columnNamesCustomerEvent},
		{"csv/Content.csv", "Content", columnNamesContent},
		{"csv/ContentPrice.csv", "ContentPrice", columnNamesContentPrice},
		{"csv/CustomerEventData.csv", "CustomerEventData", columNamesCustomerEventData},
	}

	for _, step := range steps {
		if err := importDataFromCSV(db, step.csvPath, step.tableName, step.columnNames); err != nil {
			log.Printf("Error %s when importing data from CSV %s\n", err, step.tableName)
			return err
		}
		// Update and render the main progress bar
		progressBar.Add(1)
	}
	progressBar.Clear()
	progressBar.Finish()

	return nil
}

func initializeProgressBar(message string) *progressbar.ProgressBar {
	Bar := progressbar.NewOptions(1000,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription(message),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
	Bar.RenderBlank()
	time.Sleep(600 * time.Millisecond)
	return Bar
}

func main() {
	var login string = getDotEnvVar("LOGIN")
	var password string = getDotEnvVar("PASSWORD")
	var ip string = getDotEnvVar("IP")
	var dbname string = "quantic_db"

	dbCreationBar := initializeProgressBar("Creating database [...]")
	dbCreationBar.RenderBlank()

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

	// Finish the database creation progress bar
	dbCreationBar.Clear()
	dbCreationBar.Finish()

	dataPopulationBar := initializeProgressBar("Filling Data with .csv files [...]")
	dataPopulationBar.RenderBlank()

	if err := createTypesTables(quanticDB); err != nil {
		// log.Printf("Error %s when creating types tables\n", err)
		return
	}

	// Import data into tables
	if err := fillAllTables(quanticDB, dataPopulationBar); err != nil {
		log.Printf("Error %s when importing data into tables\n", err)
		return
	}

}
