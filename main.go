package main

/*
Import github.com/go-sql-driver/mysql
Permit to import MySQL
*/
import (
	// "context"
	// "encoding/csv"
	// "fmt"
	"log"
	"sort"

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
// type Customer struct {
// 	CustomerID       int64
// 	ClientCustomerID int64
// 	InsertDate       time.Time
// }

type CustomerSpent struct {
	CustomerID int64
	Spent      float64
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
	startDate := time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC)
	customersMoneySpent, err := getCustomerSpentMap(quanticDB, startDate)
	if err != nil {
		return
	}
	printRandomEntriesMap(customersMoneySpent)

	// Convert map values to a slice
	// Create a slice of CustomerSpent
	moneySpentSlice := make([]CustomerSpent, 0, len(customersMoneySpent))
	for customerID, spent := range customersMoneySpent {
		moneySpentSlice = append(moneySpentSlice, CustomerSpent{CustomerID: customerID, Spent: spent})
	}

	// Sort the slice by Spent in descending order
	sort.Slice(moneySpentSlice, func(i, j int) bool {
		return moneySpentSlice[i].Spent > moneySpentSlice[j].Spent
	})

	// Iterate over the sorted slice
	// for _, entry := range moneySpentSlice {
	// 	fmt.Printf("CustomerID: %d, Spent: %.2f\n", entry.CustomerID, entry.Spent)
	// }
}
