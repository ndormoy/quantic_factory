package main

import (
    "database/sql"
    "fmt"
    "log"
    "time"

    _ "github.com/go-sql-driver/mysql"
)

// Function to save data to the 'test_export_DATE' table
func saveDataToTable(db *sql.DB, date time.Time, data []CustomerData) error {
    // Define the table name based on the date (YYYYMMDD format)
    tableName := date.Format("20060102")

    // Check if the table exists
    tableExists, err := doesTableExist(db, tableName)
    if err != nil {
        return err
    }

    // If the table doesn't exist, create it
    if !tableExists {
        err = createTable(db, tableName)
        if err != nil {
            return err
        }
    }

    // Iterate through the data and update or insert rows
    for _, entry := range data {
        if entryExists(db, tableName, entry.CustomerID) {
            err = updateRow(db, tableName, entry)
            if err != nil {
                return err
            }
        } else {
            err = insertRow(db, tableName, entry)
            if err != nil {
                return err
            }
        }
    }

    return nil
}

// Function to check if a table exists
func doesTableExist(db *sql.DB, tableName string) (bool, error) {
    query := fmt.Sprintf("SHOW TABLES LIKE '%s'", tableName)
    rows, err := db.Query(query)
    if err != nil {
        return false, err
    }
    defer rows.Close()

    return rows.Next(), nil
}

// Function to create a table
func createTable(db *sql.DB, tableName string) error {
    query := fmt.Sprintf(`CREATE TABLE %s (
        CustomerID INT PRIMARY KEY,
        Email VARCHAR(255),
        CA DECIMAL(10, 2)
    )`, tableName)

    _, err := db.Exec(query)
    return err
}

// Function to check if a CustomerID exists in a table
func entryExists(db *sql.DB, tableName string, customerID int64) bool {
    query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE CustomerID = ?", tableName)
    var count int
    err := db.QueryRow(query, customerID).Scan(&count)
    if err != nil {
        log.Println(err)
    }
    return count > 0
}

// Function to update a row in the table
func updateRow(db *sql.DB, tableName string, data CustomerData) error {
    query := fmt.Sprintf("UPDATE %s SET CA = ? WHERE CustomerID = ?", tableName)
    _, err := db.Exec(query, data.CA, data.CustomerID)
    return err
}

// Function to insert a new row into the table
func insertRow(db *sql.DB, tableName string, data CustomerData) error {
    query := fmt.Sprintf("INSERT INTO %s (CustomerID, Email, CA) VALUES (?, ?, ?)", tableName)
    _, err := db.Exec(query, data.CustomerID, data.Email, data.CA)
    return err
}

// Define the structure for customer data
type CustomerData struct {
    CustomerID int64
    Email      string
    CA         float64
}

// func main() {
//     // Establish a database connection and obtain a 'db' handle

//     // Define the date and customer data
//     date := time.Now()
//     data := []CustomerData{
//         {CustomerID: 1, Email: "customer1@example.com", CA: 1000.50},
//         {CustomerID: 2, Email: "customer2@example.com", CA: 750.25},
//     }

//     // Call the saveDataToTable function to save the data
//     err := saveDataToTable(db, date, data)
//     if err != nil {
//         log.Fatal(err)
//     }

//     // Close the database connection when done
//     defer db.Close()
// }


func manageExport(db *sql.DB, moneySpentSlice []CustomerSpent) {
    // Establish a database connection and obtain a 'db' handle

    // Define the date and customer data
    date := time.Now()
    data := []CustomerData{
        {CustomerID: 1, Email: "customer1@example.com", CA: 1000.50},
        {CustomerID: 2, Email: "customer2@example.com", CA: 750.25},
    }

    // Call the saveDataToTable function to save the data
    err := saveDataToTable(db, date, data)
    if err != nil {
        log.Fatal(err)
    }

    // Close the database connection when done
    defer db.Close()
}