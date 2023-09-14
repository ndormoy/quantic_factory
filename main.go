package main

/*
Import github.com/go-sql-driver/mysql
Permit to import MySQL
*/
import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type CustomerSpent struct {
	CustomerID int64
	Spent      float64
}

type structCustomerEventData struct {
	EventID     int64
	ContentID   int64
	CustomerID  int64
	EventTypeID int64
	EventDate   string
	Quantity    int
	InsertDate  string
}

func main() {
	var login string = getDotEnvVar("LOGIN")
	var password string = getDotEnvVar("PASSWORD")
	var ip string = getDotEnvVar("IP")
	var dbname string = "quantic_db"
	var quanticDB *sql.DB
	var startDate = time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC)

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

	//----------------------------------------------------------------------

	currentDate := time.Now().Format("20060102")

	/*
		Uncomment to test if the Export Tables change when some CustomerEvent is added
		with eventType 6, on customer 2065908675
	*/
	// err = createCustomerEventData(quanticDB, structCustomerEventData{
	// 	EventID:     1656516851, // Specify the actual values here
	// 	ContentID:   1726958166,
	// 	CustomerID:  2065908675,
	// 	EventTypeID: 6,
	// 	EventDate:   currentDate, // Use the appropriate date
	// 	Quantity:    1,           // Specify the quantity
	// 	InsertDate:  currentDate, // Use the appropriate date
	// })
	// if err != nil {
	// 	log.Printf("%s\n", err)
	// 	return
	// }
	// log.Printf("New CustomerEventData created\n")
	//----------------------------------------------------------------------

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
	sortMySlices(&spentValues, &moneySpentSlice)

	_, err = createBestClientMap(spentValues, moneySpentSlice)
	if err != nil {
		log.Printf("Error when creating the map with the best clients: %s\n", err)
		return
	}

	printSpentSlice(moneySpentSlice)

	CalculateQuantilesNearestRank(spentValues, 40)

	manageExport(quanticDB, moneySpentSlice, currentDate)
}
