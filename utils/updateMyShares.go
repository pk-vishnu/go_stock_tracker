package utils

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/hashicorp/go-multierror"
)

// Limit for concurrent goroutines to prevent resource exhaustion
const maxConcurrency = 20

func UpdateMyShares() {
	db, err := dbConnect()
	if err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
	}
	defer db.Close()

	// Fetch tickers in one batch
	rows, err := db.Query("SELECT ticker FROM MyShares")
	if err != nil {
		log.Fatalf("Error fetching tickers: %v", err)
	}
	defer rows.Close()

	// Collect tickers
	var tickers []string
	for rows.Next() {
		var ticker string
		if err := rows.Scan(&ticker); err != nil {
			log.Printf("Error scanning ticker: %v", err)
			continue
		}
		tickers = append(tickers, ticker)
	}

	// Fetch instrument keys in one query to avoid redundant queries
	instrumentKeys := make(map[string]string)
	for _, ticker := range tickers {
		var instrumentKey string
		err := db.QueryRow("SELECT instrument_key FROM TickerKeys WHERE ticker = ?", ticker).Scan(&instrumentKey)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("No instrument key found for ticker: %s", ticker)
			} else {
				log.Printf("Error fetching instrument key for ticker %s: %v", ticker, err)
			}
			continue
		}
		instrumentKeys[ticker] = instrumentKey
	}

	// Use a semaphore to limit concurrency
	sem := make(chan struct{}, maxConcurrency)

	var wg sync.WaitGroup
	var resultErr error
	errChan := make(chan error, len(tickers))

	for _, ticker := range tickers {
		wg.Add(1)
		// Limit concurrency
		sem <- struct{}{}

		go func(ticker string) {
			defer func() {
				wg.Done()
				<-sem
			}()

			instrumentKey, exists := instrumentKeys[ticker]
			if !exists {
				log.Printf("No instrument key found for ticker: %s", ticker)
				return
			}

			stockData, err := GetLiveStockData(ticker)
			if err != nil {
				log.Printf("Error fetching live stock data for ticker %s: %v", ticker, err)
				errChan <- fmt.Errorf("error fetching live stock data for %s: %v", ticker, err)
				return
			}

			sma20 := CalculateSMA20(instrumentKey)
			sma50 := CalculateSMA50(instrumentKey)

			_, err = db.Exec(`
				UPDATE MyShares
				SET current_price = ?, sma20 = ?, sma50 = ?, last_updated = CURRENT_TIMESTAMP
				WHERE ticker = ?`,
				stockData.Price, sma20, sma50, ticker)
			if err != nil {
				log.Printf("Error updating MyShares for ticker %s: %v", ticker, err)
				errChan <- fmt.Errorf("error updating MyShares for %s: %v", ticker, err)
			} else {
				log.Printf("Successfully updated ticker %s", ticker)
			}
		}(ticker)
	}


	wg.Wait()
	close(errChan)

	// Aggregate all errors
	for err := range errChan {
		resultErr = multierror.Append(resultErr, err)
	}

	if resultErr != nil {
		log.Printf("Errors encountered: %v", resultErr)
	}
}
