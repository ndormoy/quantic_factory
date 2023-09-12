package main

/*
Import github.com/go-sql-driver/mysql
Permit to import MySQL
*/
import (
	// "context"
	// "encoding/csv"
	"fmt"
	"log"
	// "os"
	// "strings"
	"time"
	// "bufio"
	// Import godotenv
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	// "github.com/joho/godotenv"
	// "github.com/schollz/progressbar/v3"
)

// Create a struct to represent Customer data
type Customer struct {
	CustomerID int64
	ClientCustomerID int64
	InsertDate time.Time
}

func main() {
	var login string = getDotEnvVar("LOGIN")
	var password string = getDotEnvVar("PASSWORD")
	var ip string = getDotEnvVar("IP")
	var dbname string = "quantic_db"
	var quanticDB *sql.DB
	// var customerRevenueMap map[int64]float64

	quanticDB, err := isQuanticDBExists(login, password, ip, dbname)
	if err != nil {
		log.Printf("Error checking if the database exists: %s\n", err)
		return
	}
	if quanticDB == nil {
		if quanticDB, err = initAndFill(login, password, ip, dbname); err != nil {
			log.Printf("Error when initializing and filling the database: %s\n", err)
			return
		}
	}




	// Initialize the customerRevenueMap
	// customerRevenueMap = make(map[int64]float64)

	// Query to retrieve CustomerEventData records with EventTypeID = 6 and ContentID
	query := `
		SELECT c.CustomerID, ce.ContentID
		FROM CustomerEventData ce
		INNER JOIN Customer c ON ce.CustomerID = c.CustomerID
		WHERE ce.EventTypeID = 6
	`

	rows, err := quanticDB.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Create a map to store CustomerID and total Price
	customerPriceMap := make(map[int64]float64)

	// Iterate through the rows to get ContentID and Price
	for rows.Next() {
		var customerID, contentID int64
		err := rows.Scan(&customerID, &contentID)
		if err != nil {
			log.Fatal(err)
		}

		// Query to retrieve Price from ContentPrice
		priceQuery := "SELECT Price FROM ContentPrice WHERE ContentID = ?"
		var price float64
		err = quanticDB.QueryRow(priceQuery, contentID).Scan(&price)
		if err != nil {
			log.Fatal(err)
		}

		// Add the Price to the customer's total Price
		customerPriceMap[customerID] += price
	}

	// Print the map with CustomerID and total Price
	for customerID, totalPrice := range customerPriceMap {
		fmt.Printf("CustomerID: %d, Total Price: %.2f\n", customerID, totalPrice)
	}
}
