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
	"sort"

	// "os"
	// "strings"
	"time"
	// "bufio"
	// Import godotenv
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	// "gonum.org/v1/gonum/stat"
	// "github.com/joho/godotenv"
	// "github.com/schollz/progressbar/v3"
)

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

	// Convert map values to a slice, Create a slice of CustomerSpent
	moneySpentSlice := make([]CustomerSpent, 0, len(customersMoneySpent))
	for customerID, spent := range customersMoneySpent {
		moneySpentSlice = append(moneySpentSlice, CustomerSpent{CustomerID: customerID, Spent: spent})
	}

	// Create a slice to hold the Spent values
	spentValues := make([]float64, len(moneySpentSlice))
	for i, entry := range moneySpentSlice {
		spentValues[i] = entry.Spent
	}

	// Sort the Spent values slice in ascending order
	sort.Float64s(spentValues)

	// Sort the moneySpentSlice by Spent in descending order
	sort.Slice(moneySpentSlice, func(i, j int) bool {
		return moneySpentSlice[i].Spent < moneySpentSlice[j].Spent
	})

	/*-------------------------------*/
	fmt.Printf("------------------------------------------------------------\n")

	_, err = createBestClientMap(spentValues, moneySpentSlice)
	if err != nil {
		log.Printf("Error when creating the map with the best clients: %s\n", err)
		return
	}

	fmt.Printf("------------------------------------------------------------\n")

	/*-------------------------------*/

	// Iterate over the sorted slice
	for _, entry := range moneySpentSlice {
		fmt.Printf("CustomerID: %d, Spent: %.2f\n", entry.CustomerID, entry.Spent)
	}

	// CreateAllQuantileMap(spentValues), 40 is the number of quantiles
	quantileInfoMap, quantiles := CalculateQuantilesNearestRank(spentValues, 40)
	for _, q := range quantiles {
		// fmt.Print("q = %f ", q)
		// fmt.Printf("%.2f%% Quantile Info:\n", q*100)
		info := quantileInfoMap[q]
		fmt.Printf("%.2f%% Quantile Info:\n", q*100 - 2.5)
		fmt.Printf("Number of clients: %d\n", info.NumClients)
		fmt.Printf("Max Revenue: %.2f\n", info.MaxRevenue)
		fmt.Printf("Min Revenue: %.2f\n\n", info.MinRevenue)
	}
	

}
