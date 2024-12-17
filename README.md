# GoLang NSE Stock Tracker

This project is a stock tracker and report generator that fetches data from a publicly available Google Sheets CSV, syncs the data to a MySQL database, updates share data (including current price, SMA20, and SMA50), and generates a report of gainers and losers. This report is sent to a Telegram bot based on user-defined thresholds and SMAs.

## Features

- **Sync Google Sheets to MySQL**: Fetch data from a Google Sheets CSV (publicly accessible) and sync it to a MySQL database. This includes data such as ticker symbols, number of shares, buying price, and threshold percentage.
- **Stock Data Updates**: Updates the current stock price, SMA20, and SMA50 for each ticker symbol.
- **Report Generation**: Automatically generates a daily report of stock gainers and losers based on user-defined thresholds and SMAs.
- **Telegram Integration**: Sends the generated report to a specified Telegram bot that alerts the user of the latest gains and losses.
