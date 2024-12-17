package utils

import (
	"fmt"
	"time"
)

func CalculateSMA20(instrumentKey string) float64 {
	var sum float64 = 0
	data, err := fetchHistoricalData(instrumentKey, time.Now().Format("2006-01-02"), time.Now().AddDate(0, 0, -20).Format("2006-01-02"))
	if err != nil {
		fmt.Println(err)
	}
	candles, err := parseHistoricalData(data)
	if err != nil {
		fmt.Println(err)
	}
	for _, candle := range candles {
		sum += candle.Close
	}
	return sum / 20
}


func CalculateSMA50(instrumentKey string) float64 {
	var sum float64 = 0
	data, err := fetchHistoricalData(instrumentKey, time.Now().Format("2006-01-02"), time.Now().AddDate(0, 0, -50).Format("2006-01-02"))
	if err != nil {
		fmt.Println(err)
	}
	candles, err := parseHistoricalData(data)
	if err != nil {
		fmt.Println(err)
	}
	for _, candle := range candles {
		sum += candle.Close
	}
	return sum / 50
}
