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
	CustomerID       int64
	ClientCustomerID int64
	InsertDate       time.Time
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

	contentIDs, err := getContentIDFromPurchase(quanticDB)
	if err != nil {
		log.Printf("Error getting ContentID in CustomerEventData where EventTypeID == 6 in function getContentIDFromPurchase : %s\n", err)
		return
	}
	currencyMap, err := getCurrencyForCustomers(quanticDB, contentIDs)
	if err != nil {
		log.Printf("Error getting Currency for Customers in function getCurrencyForCustomers : %s\n", err)
		return
	}
	fmt.Printf("%v", currencyMap)
	fmt.Printf("--------------------------------\n")
	customerIDs, err := calculateTotalPurchaseAmounts(quanticDB, contentIDs, currencyMap)
	if err != nil {
		log.Printf("Error when creating and return a map with CustomerID and their purchases, in function createMapWithCustomerIDPurchase : %s\n", err)
		return
	}
	// fmt.Printf("customerIDs : %v\n", customerIDs)
	// Print the map with CustomerID and total Purchase Amount
	for customerID, purchaseAmount := range customerIDs {
		fmt.Printf("CustomerID: %d, Total Purchase Amount: %.2f\n", customerID, purchaseAmount)
	}
}
