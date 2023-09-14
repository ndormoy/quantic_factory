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

// func createNewCustomerEventData(db *sql.DB, customerID int64, eventTypeID int64, eventID int64, contentID int64) error {
// 	// Generate the date in YYYYMMDD format
// 	currentDate := time.Now().Format("20060102")

// 	// Define the event type data
// 	eventData := struct {
// 		CustomerID  int64
// 		EventTypeID int64
// 		EventID     int64
// 		EventDate   string
// 		Quantity    int
// 		ContentID   int64
// 	}{
// 		CustomerID:  customerID,
// 		EventTypeID: eventTypeID,
// 		EventID:     eventID,
// 		EventDate:   currentDate,
// 		ContentID:   contentID,
// 		Quantity:    1,
// 	}

// 	// Create new customer event
// 	fmt.Printf("DEBUG : %d\n", eventData.EventID)
// 	if err := createNewCustomerEvent(db, eventData, currentDate); err != nil {
// 		return err
// 	}

// 	// Commit the transaction for CustomerEvent table
// 	_, err := db.Exec("COMMIT")
// 	if err != nil {
// 		fmt.Print("Failed to commit transaction for CustomerEvent: ", err)
// 		return err
// 	}

// 	// Insert the event data into the CustomerEventData table
// 	_, err = db.Exec(`
//         INSERT INTO CustomerEventData (EventID, ContentID, CustomerID, EventTypeID, EventDate, Quantity, InsertDate)
//         VALUES (?, ?, ?, ?, ?, ?, ?)
//     `, eventData.EventID, eventData.ContentID, eventData.CustomerID, eventData.EventTypeID, eventData.EventDate, eventData.Quantity, currentDate)

// 	if err != nil {
// 		fmt.Print("YO")
// 		return fmt.Errorf("error creating a new event type: %w", err)
// 	}

// 	// Commit the transaction for CustomerEventData table
// 	_, err = db.Exec("COMMIT")
// 	if err != nil {
// 		fmt.Print("Failed to commit transaction for CustomerEventData: ", err)
// 		return err
// 	}

// 	// Query to check if the new row exists in CustomerEventData
// 	query := "SELECT COUNT(*) FROM CustomerEventData WHERE EventID = ?"
// 	var count int
// 	err = db.QueryRow(query, eventData.EventID).Scan(&count)
// 	if err != nil {
// 		fmt.Print("HAHAHA")
// 		log.Printf("Error checking if the new row exists: %s\n", err)
// 	} else {
// 		if count > 0 {
// 			log.Printf("New row in CustomerEventData exists")
// 		} else {
// 			log.Printf("New row in CustomerEventData does not exist")
// 		}
// 	}

// 	return nil
// }



func createCustomerEventData(db *sql.DB, eventData structCustomerEventData) error {

	// Create new customer event
	fmt.Printf("DEBUG : %d\n", eventData.EventID)
	if err := createNewCustomerEvent(db, eventData); err != nil {
		return err
	}

	// Insert the event data into the CustomerEventData table
	_, err := db.Exec(`
        INSERT INTO CustomerEventData (EventDataID, EventID, ContentID, CustomerID, EventTypeID, EventDate, Quantity, InsertDate)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `, 2222222222, eventData.EventID, eventData.ContentID, eventData.CustomerID, eventData.EventTypeID, eventData.EventDate, eventData.Quantity, eventData.InsertDate)

	if err != nil {
		return fmt.Errorf("error creating a new CustomerEventData: %w", err)
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
