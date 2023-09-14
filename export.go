package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func manageExport(db *sql.DB, moneySpentSlice []CustomerSpent, currentDate string) error {

	// Combine with a prefix test_export_
	var tableName string = "test_export_" + currentDate

	// Activate for testing creating tables on different days
	// var tableName string = "test_export_20230915"
	// Create the 'test_export_DATE' table if it doesn't exist
	_, err := db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			DATE DATE,
			CustomerID BIGINT,
			CA DOUBLE,
			PRIMARY KEY (DATE, CustomerID)
		)
	`, tableName))
	if err != nil {
		return fmt.Errorf("error creating the table: %w", err)
	}

	// Begin a transaction for the bulk insert
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	// Prepare the bulk insert statement
	stmt, err := tx.Prepare(fmt.Sprintf("INSERT INTO %s (DATE, CustomerID, CA) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE CA = ?", tableName))
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	// Perform the bulk insert
	for _, entry := range moneySpentSlice {
		_, err := stmt.Exec(currentDate, entry.CustomerID, entry.Spent, entry.Spent)
		if err != nil {
			return fmt.Errorf("error executing statement: %w", err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

/*
	Fonction that creates a new CustomerEvent and a new CustomerEventData, to test the export
*/

func createCustomerEventData(db *sql.DB, eventData structCustomerEventData) error {
	// Create new customer event
	if err := createNewCustomerEvent(db, eventData); err != nil {
		return err
	}

	// Check if the event data with EventDataID 2222222222 already exists in CustomerEventData
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM CustomerEventData WHERE EventDataID = ?", 2222222222).Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking if the CustomerEventData entry exists: %w", err)
	}

	if count == 0 {
		// Insert the event data into the CustomerEventData table
		_, err := db.Exec(`
			INSERT INTO CustomerEventData (EventDataID, EventID, ContentID, CustomerID, EventTypeID, EventDate, Quantity, InsertDate)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, 2222222222, eventData.EventID, eventData.ContentID, eventData.CustomerID, eventData.EventTypeID, eventData.EventDate, eventData.Quantity, eventData.InsertDate)

		if err != nil {
			return fmt.Errorf("error creating a new CustomerEventData: %w", err)
		}
	}

	return nil
}

func createNewCustomerEvent(db *sql.DB, eventData structCustomerEventData) error {
	var ClientEventID int64 = 11111111111

	// Query to check if the event with the given EventID already exists
	var eventExists int
	err := db.QueryRow("SELECT COUNT(*) FROM CustomerEvent WHERE EventID = ?", eventData.EventID).Scan(&eventExists)
	if err != nil {
		return fmt.Errorf("error checking if the event exists: %w", err)
	}

	if eventExists == 0 {
		// The event does not exist, so insert it
		_, err = db.Exec(`
			INSERT INTO CustomerEvent (EventID, ClientEventID, InsertDate)
			VALUES (?, ?, ?)
		`, eventData.EventID, ClientEventID, eventData.InsertDate)

		if err != nil {
			return fmt.Errorf("error creating a new customer event: %w", err)
		}
	}

	return nil
}
