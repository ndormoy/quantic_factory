package main

import (
	// "database/sql"

	// "fmt"
	"fmt"
	"log"
	"strings"
	"time"

	// "time"
	// "context"
	// "os"
	// "encoding/csv"
	// "strings"
	"database/sql"
)

/*
Create all maps to return the map with key = CustomerID and value = moneySpent
*/

func getCustomerSpentMap(db *sql.DB, startDate time.Time) (map[int64]float64, error) {
	contentIDs, err := getContentIDFromPurchaseAfterDate(db, startDate)
	if err != nil {
		log.Printf("Error getting ContentID after %s in function getContentIDFromPurchaseAfterDate: %s\n", startDate.Format("2006-01-02"), err)
		return nil, err
	}
	currencyMap, err := getCurrencyForCustomers(db, contentIDs)
	if err != nil {
		log.Printf("Error getting Currency for Customers in function getCurrencyForCustomers: %s\n", err)
		return nil, err
	}
	customerIDs, err := calculateTotalPurchaseAmountsAfterDate(db, contentIDs, currencyMap, startDate)
	if err != nil {
		log.Printf("Error when creating and returning a map with CustomerID and their purchases in function calculateTotalPurchaseAmounts: %s\n", err)
		return nil, err
	}
	// for customerID, purchaseAmount := range customerIDs {
	// 	fmt.Printf("CustomerID: %d, Total Purchase Amount: %.2f\n", customerID, purchaseAmount)
	// }
	return customerIDs, nil
}

/*
This function get back the ContentID in CustomerEventData where EventTypeID == 6 (purchase)
We calculate the total purchase amount starting from the startDate
*/

func getContentIDFromPurchaseAfterDate(db *sql.DB, startDate time.Time) ([]int64, error) {
	query := `
    SELECT ContentID
    FROM CustomerEventData
    WHERE EventTypeID = 6 AND InsertDate >= ?
    `
	rows, err := db.Query(query, startDate)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
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
We return something like this, a map : CustomerIDs = moneySpent
*/

func calculateTotalPurchaseAmountsAfterDate(db *sql.DB, contentIDs []int64, currencyMap map[int64]string, startDate time.Time) (map[int64]float64, error) {
	customerPurchaseAmounts := make(map[int64]float64)
	processedPurchases := make(map[int64]map[int64]bool) // Track processed purchases per customer

	for _, contentID := range contentIDs {
		query := `
            SELECT CustomerID, Quantity
            FROM CustomerEventData
            WHERE ContentID = ? AND EventTypeID = 6 AND InsertDate >= ?
        `
		rows, err := db.Query(query, contentID, startDate)
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

			// Fetch the currency for the contentID
			currency, ok := currencyMap[contentID]
			if !ok {
				return nil, fmt.Errorf("Currency not found for ContentID: %d", contentID)
			}

			// Calculate the purchase amount for this specific purchase
			priceQuery := "SELECT Price FROM ContentPrice WHERE ContentID = ?"
			var price float64
			err := db.QueryRow(priceQuery, contentID).Scan(&price)
			if err != nil {
				return nil, err
			}
			purchaseAmount := price * float64(quantity)

			// Convert the purchase amount to euros
			if currency != "EUR" {
				convertedAmount, err := convertToEUR(purchaseAmount, currency)
				if err != nil {
					return nil, err
				}
				purchaseAmount = convertedAmount
			}
			// Add the purchase amount to the customer's total
			customerPurchaseAmounts[customerID] += purchaseAmount
			processedPurchases[customerID][contentID] = true
		}
	}
	return customerPurchaseAmounts, nil
}

/*
This function return a map containing : key = ContentID, value = Currency.
Permit later to convert the purchase amount to Euros.
*/

func getCurrencyForCustomers(db *sql.DB, contentIDs []int64) (map[int64]string, error) {
	currencyMap := make(map[int64]string)

	// Query to get Currency for multiple ContentIDs using IN clause
	query := `
        SELECT ContentID, Currency
        FROM ContentPrice
        WHERE ContentID IN (?)
    `

	// Create a comma-separated list of contentIDs for the query
	contentIDStr := ""
	for i, contentID := range contentIDs {
		contentIDStr += fmt.Sprintf("%d", contentID)
		if i < len(contentIDs)-1 {
			contentIDStr += ","
		}
	}

	// Replace the placeholder with the contentID list in the query
	query = strings.Replace(query, "?", contentIDStr, 1)

	rows, err := db.Query(query)
	if err != nil {
		fmt.Printf("POUET")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var contentID int64
		var currency string
		if err := rows.Scan(&contentID, &currency); err != nil {
			fmt.Printf("zeubi")
			return nil, err
		}
		currencyMap[contentID] = currency
	}
	return currencyMap, nil
}
