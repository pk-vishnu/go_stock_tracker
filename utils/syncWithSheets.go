package utils

import (
	"encoding/csv"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// SyncGoogleSheetToDB fetches data from a public Google Sheets CSV and updates the MySQL table, deleting records no longer in the sheet.
func SyncGoogleSheetToDB() {

	// Connect to the database
	db, err := dbConnect()
	if err != nil {
		log.Fatal("Error connecting to DB")
	}
	defer db.Close()

	// Google Sheets public CSV URL (replace with actual document ID and sheet name)
	csvURL := "https://docs.google.com/spreadsheets/d/1uXx5P3I9vCjesAOcxcd_KAL68dLAs856fnEsjPuxM_E/gviz/tq?tqx=out:csv&sheet=golang_stock_screener"

	// Fetch the CSV file from the Google Sheets public URL
	resp, err := http.Get(csvURL)
	if err != nil {
		log.Fatalf("Failed to fetch CSV: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to fetch CSV: status code %d", resp.StatusCode)
	}

	// Parse the CSV file
	reader := csv.NewReader(resp.Body)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to parse CSV: %v", err)
	}

	// Create a map to track tickers in the Google Sheets
	googleSheetTickers := make(map[string]bool)

	// Process CSV rows (skip the header)
	for i, row := range records {
		if i == 0 || len(row) < 4 { // Skip the header or malformed rows
			continue
		}

		// Extract data from CSV
		ticker := strings.TrimSpace(row[0])
		numShares, _ := strconv.Atoi(strings.TrimSpace(row[1]))
		buyingPrice, _ := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
		threshold, _ := strconv.ParseFloat(strings.TrimSpace(row[3]), 64)

		// Mark the ticker as present in the sheet
		googleSheetTickers[ticker] = true

		// Sync MySQL table: Insert or update record
		_, err := db.Exec(`
			INSERT INTO MyShares (ticker, num_shares, buying_price, threshold)
			VALUES (?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE 
			num_shares=VALUES(num_shares), 
			buying_price=VALUES(buying_price), 
			threshold=VALUES(threshold)
		`, ticker, numShares, buyingPrice, threshold)

		if err != nil {
			log.Printf("Failed to update ticker %s: %v", ticker, err)
		} else {
			log.Printf("Updated ticker %s successfully", ticker)
		}
	}

	// Delete tickers that are no longer in the Google Sheets CSV
	deleteSQL := "DELETE FROM MyShares WHERE ticker NOT IN ("
	args := []interface{}{}
	for ticker := range googleSheetTickers {
		deleteSQL += "?,"
		args = append(args, ticker)
	}

	// Remove the trailing comma
	deleteSQL = deleteSQL[:len(deleteSQL)-1] + ")"

	// Execute the delete query
	_, err = db.Exec(deleteSQL, args...)
	if err != nil {
		log.Printf("Error deleting obsolete records: %v", err)
	} else {
		log.Println("Deleted obsolete records successfully")
	}

	log.Println("Sync complete!")
}
