package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

//Function to get Live Stock Data by scraping google finance
type StockData struct {
	Price   float64
	PERatio float64
}

func GetLiveStockData(stockName string) (StockData, error) {
	url := fmt.Sprintf("https://www.google.com/finance/quote/%s:NSE", stockName)
	c := colly.NewCollector()

	var stockData StockData
	var priceStr, peRatioStr string

	//current price
	c.OnHTML("div.YMlKec.fxKbKc", func(e *colly.HTMLElement) {
		priceStr = strings.ReplaceAll(e.Text, "₹", "") 
		priceStr = strings.ReplaceAll(priceStr, ",", "")
	})

	//P/E Ratio
	c.OnHTML("div.eYanAe", func(e *colly.HTMLElement) {
		e.ForEach("div.gyFHrc", func(index int, el *colly.HTMLElement) {
			if index == 5 {
				peRatioStr = el.ChildText("div.P6K39c")
				if peRatioStr != "" {
					peRatioStr = strings.ReplaceAll(peRatioStr, ",", "")
				}
			}
		})
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error: %v, Status Code: %d", err, r.StatusCode)
	})

	err := c.Visit(url)
	if err != nil {
		return stockData, err
	}

	//Parsing data into floats
	if priceStr != "" {
		stockData.Price, err = strconv.ParseFloat(priceStr, 64)
		if err != nil {
			return stockData, fmt.Errorf("failed to parse price: %v", err)
		}
	}
	if peRatioStr != "" {
		stockData.PERatio, err = strconv.ParseFloat(peRatioStr, 64)
		if err != nil {
			return stockData, fmt.Errorf("failed to parse P/E ratio: %v", err)
		}
	}

	return stockData, nil
}

//Functions to fetch and parse historical Data from UPSTOX 
func fetchHistoricalData(instrumentKey, fromDate, toDate string)([]byte, error){
	url := fmt.Sprintf("https://api.upstox.com/v2/historical-candle/%s/day/%s/%s/",instrumentKey, fromDate, toDate)
	req, err := http.NewRequest("GET", url, nil)
	if err!=nil{
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err!=nil{
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err!=nil{
		return nil, err
	}
	return body, nil
}

type HistoricalData struct{
	Date string `json:"date"`
	Close float64 `json:"close"`
}
type APIResponse struct{
	Data struct{
		Candles [][]interface{} `json:"candles"`
	}`json:"data"`
}

func parseHistoricalData(data []byte)([]HistoricalData, error){
	var apiResponse APIResponse
	if err := json.Unmarshal(data, &apiResponse); err!=nil{
		return nil, err
	}
	
	var candles []HistoricalData 
	for _, item := range apiResponse.Data.Candles{
		if len(item)>4{
			date, close := item[0].(string), item[4].(float64)
			candles = append(candles, HistoricalData{Date: date, Close: close})
		}
	}
	return candles, nil
}

