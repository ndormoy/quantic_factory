package main

import (
	// "database/sql"
	"fmt"
	// "fmt"
	"log"

	// "time"
	// "context"
	// "os"
	// "encoding/csv"
	// "strings"
	"database/sql"
)

func mapWithCustIDPurchaseJOIN(db *sql.DB) {
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

/*
This function get back the ContentID in CustomerEventData where EventTypeID == 6 (purchase)
*/

func getContentIDFromPurchase(db *sql.DB) ([]int64, error) {
	query := `
	SELECT ContentID
	FROM CustomerEventData
	WHERE EventTypeID = 6
	`
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	// to ensure proper resource cleanup
	defer rows.Close()

	var contentIDs []int64

	for rows.Next() {
		var contentID int64
		if err := rows.Scan(&contentID); err != nil {
			log.Fatal(err)
			return nil, err
		}
		contentIDs = append(contentIDs, contentID)
	}
	return contentIDs, nil
}

/*
This function create and return a map with CustomerID and their purchases we multiplie with Quantity of a product.
We use ContentID from the function getContentIDFromPurchase to retrieve the Price from ContentPrice.
*/

func calculateTotalPurchaseAmounts(db *sql.DB, contentIDs []int64) (map[int64]float64, error) {
	customerPurchaseAmounts := make(map[int64]float64)
	processedPurchases := make(map[int64]map[int64]bool) // Track processed purchases per customer

	for _, contentID := range contentIDs {
		query := `
            SELECT CustomerID, Quantity
            FROM CustomerEventData
            WHERE ContentID = ? AND EventTypeID = 6
        `
		rows, err := db.Query(query, contentID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var customerID, quantity int64
			if err := rows.Scan(&customerID, &quantity); err != nil {
				return nil, err
			}

			// Check if this purchase has already been processed for this customer
			if processedPurchases[customerID] == nil {
				processedPurchases[customerID] = make(map[int64]bool)
			}
			if processedPurchases[customerID][contentID] {
				continue // Skip if already processed
			}

			// Calculate the purchase amount for this specific purchase
			priceQuery := "SELECT Price FROM ContentPrice WHERE ContentID = ?"
			var price float64
			err := db.QueryRow(priceQuery, contentID).Scan(&price)
			if err != nil {
				return nil, err
			}
			purchaseAmount := price * float64(quantity)

			// Add the purchase amount to the customer's total
			customerPurchaseAmounts[customerID] += purchaseAmount
			processedPurchases[customerID][contentID] = true
			// fmt.Printf("Added purchase amount %.2f to CustomerID %d\n", purchaseAmount, customerID)
		}
	}
	return customerPurchaseAmounts, nil
}
