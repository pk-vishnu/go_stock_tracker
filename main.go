package main

import (
	"fmt"
	"log"
	"stock-screener/utils"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Stock Screener v1.0")
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	utils.SyncGoogleSheetToDB()
	utils.UpdateMyShares()
	utils.GenerateReport()
}