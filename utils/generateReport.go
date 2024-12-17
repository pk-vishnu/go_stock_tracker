package utils

import (
	"bytes"
	"database/sql"
	"html/template"
	"log"
	"sort"
	"time"
)

type MyShares struct {
    Ticker       string
    NumShares    int
    BuyingPrice  float64
    CurrentPrice sql.NullFloat64 
    Sma20        sql.NullFloat64 
    Sma50        sql.NullFloat64
    Threshold    float64
    Profit       float64
}
type ReportData struct {
	PriceGainers     []MyShares
	SMA20Gainers     []MyShares
	SMA20Losers      []MyShares
	SMA50Gainers     []MyShares
	SMA50Losers      []MyShares
	ReportTime string
}

func GenerateReport() {
	db, err := dbConnect()
	if err != nil {
		log.Fatal("Error connecting to DB")
	}
	defer db.Close()

	rows, err := db.Query("SELECT ticker, num_shares, buying_price, current_price, sma20, sma50, threshold from MyShares")
	if err != nil {
		log.Fatal("Query to MyShares failed")
	}
	defer rows.Close()

	var priceGainers []MyShares
	var sma20Gainers []MyShares
	var sma20Losers []MyShares
	var sma50Gainers []MyShares
	var sma50Losers []MyShares
	currentTime := time.Now().Format("02-Jan-2006 15:04:05")
	// Loop through all the rows in MyShares
	for rows.Next() {
		var stock MyShares
		err := rows.Scan(&stock.Ticker, &stock.NumShares, &stock.BuyingPrice, &stock.CurrentPrice, &stock.Sma20, &stock.Sma50, &stock.Threshold)
		if err != nil {
			log.Fatal("Error scanning row: ", err)
		}
		if !stock.CurrentPrice.Valid || !stock.Sma20.Valid || !stock.Sma50.Valid || stock.CurrentPrice.Float64==0.0{
			log.Printf("Skipping stock %s due to NULL value in essential field", stock.Ticker)
			continue
		}

		stock.Profit = (stock.CurrentPrice.Float64 - stock.BuyingPrice) * float64(stock.NumShares)

		if stock.CurrentPrice.Float64 > (stock.BuyingPrice * (1 + stock.Threshold/100)) {
			priceGainers = append(priceGainers, stock)
		}

		if stock.CurrentPrice.Float64 > stock.Sma20.Float64 {
			sma20Gainers = append(sma20Gainers, stock)
		} else {
			sma20Losers = append(sma20Losers, stock)
		}

		if stock.CurrentPrice.Float64 > stock.Sma50.Float64 {
			sma50Gainers = append(sma50Gainers, stock)
		} else {
			sma50Losers = append(sma50Losers, stock)
		}

	}

	sort.Slice(sma20Gainers, func(i, j int) bool {
		return sma20Gainers[i].CurrentPrice.Float64-sma20Gainers[i].Sma20.Float64 > sma20Gainers[j].CurrentPrice.Float64-sma20Gainers[j].Sma20.Float64
	})
	sort.Slice(sma20Losers, func(i, j int) bool {
		return sma20Losers[i].Sma20.Float64-sma20Losers[i].CurrentPrice.Float64 > sma20Losers[j].Sma20.Float64-sma20Losers[j].CurrentPrice.Float64
	})

	sort.Slice(sma50Gainers, func(i, j int) bool {
		return sma50Gainers[i].CurrentPrice.Float64-sma50Gainers[i].Sma50.Float64 > sma50Gainers[j].CurrentPrice.Float64-sma50Gainers[j].Sma50.Float64
	})
	sort.Slice(sma50Losers, func(i, j int) bool {
		return sma50Losers[i].Sma50.Float64-sma50Losers[i].CurrentPrice.Float64 > sma50Losers[j].Sma50.Float64-sma50Losers[j].CurrentPrice.Float64
	})

	if len(sma20Gainers) > 5 {
		sma20Gainers = sma20Gainers[:5]
	}
	if len(sma20Losers) > 5 {
		sma20Losers = sma20Losers[:5]
	}
	if len(sma50Gainers) > 5 {
		sma50Gainers = sma50Gainers[:5]
	}
	if len(sma50Losers) > 5 {
		sma50Losers = sma50Losers[:5]
	}

	reportData := ReportData{
		PriceGainers: priceGainers,
		SMA20Gainers: sma20Gainers,
		SMA20Losers:  sma20Losers,
		SMA50Gainers: sma50Gainers,
		SMA50Losers:  sma50Losers,
		ReportTime:   currentTime,
	}

	// Parse and execute the report template
	tmpl := reportTemplate()
	t, err := template.New("report").Parse(tmpl)
	if err != nil {
		log.Fatal("Error parsing report template: ", err)
	}

	var htmlReport bytes.Buffer
	err = t.Execute(&htmlReport, reportData)
	if err != nil {
		log.Fatal("Error executing template: ", err)
	}

	if err != nil {
		log.Fatal("Error writing to file: ", err)
	}

	log.Println("Report generated successfully")
	sendToTelegram(htmlReport.String())
}
