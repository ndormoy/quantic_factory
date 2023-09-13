package main

import (
	// "database/sql"
	"fmt"
	"log"
	// "time"
	// "context"
	// "os"
	// "encoding/csv"
	// "strings"
	"database/sql"

)

func mapWithCustIDPurchaseJOIN (db *sql.DB){
	// Query to retrieve CustomerEventData records with EventTypeID = 6 and ContentID
	query := `
		SELECT c.CustomerID, ce.ContentID
		FROM CustomerEventData ce
		INNER JOIN Customer c ON ce.CustomerID = c.CustomerID
		WHERE ce.EventTypeID = 6
	`

	rows, err := db.Query(query)
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
		err = db.QueryRow(priceQuery, contentID).Scan(&price)
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


