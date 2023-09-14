package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func manageExport(db *sql.DB, moneySpentSlice []CustomerSpent) error {

	// Create the 'test_export_DATE' table if it doesn't exist
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS test_export_DATE (
            DATE DATE,
            CustomerID BIGINT,
            CA DOUBLE,
            PRIMARY KEY (DATE, CustomerID)
        )
    `)
	if err != nil {
		return fmt.Errorf("error creating the table: %w", err)
	}

	// Generate the date in YYYYMMDD format
	currentDate := time.Now().Format("20060102")

	// Begin a transaction for the bulk insert
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	// Prepare the bulk insert statement
	stmt, err := tx.Prepare("INSERT INTO test_export_DATE (DATE, CustomerID, CA) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE CA = ?")
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


// func createCustomerEventData(db *sql.DB, eventData structCustomerEventData) error {

// 	// Create new customer event
// 	fmt.Printf("DEBUG : %d\n", eventData.EventID)
// 	if err := createNewCustomerEvent(db, eventData); err != nil {
// 		return err
// 	}

// 	// Insert the event data into the CustomerEventData table
// 	_, err := db.Exec(`
//         INSERT INTO CustomerEventData (EventDataID, EventID, ContentID, CustomerID, EventTypeID, EventDate, Quantity, InsertDate)
//         VALUES (?, ?, ?, ?, ?, ?, ?, ?)
//     `, 2222222222, eventData.EventID, eventData.ContentID, eventData.CustomerID, eventData.EventTypeID, eventData.EventDate, eventData.Quantity, eventData.InsertDate)

// 	if err != nil {
// 		return fmt.Errorf("error creating a new CustomerEventData: %w", err)
// 	}

// 	return nil
// }

func createCustomerEventData(db *sql.DB, eventData structCustomerEventData) error {
    // Create new customer event
    fmt.Printf("DEBUG : %d\n", eventData.EventID)
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
		fmt.Print("MDR")
		return fmt.Errorf("error checking if the event exists: %w", err)
	}

	if eventExists == 0 {
		// The event does not exist, so insert it
		_, err = db.Exec(`
            INSERT INTO CustomerEvent (EventID, ClientEventID, InsertDate)
            VALUES (?, ?, ?)
        `, eventData.EventID, ClientEventID, eventData.InsertDate)

		if err != nil {
			fmt.Print("MDR")
			return fmt.Errorf("error creating a new customer event: %w", err)
		}
	}

	return nil
}
