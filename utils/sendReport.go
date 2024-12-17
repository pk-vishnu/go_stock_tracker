package utils

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

func sendToTelegram(htmlContent string) error {
	botToken := os.Getenv("BOT_TOKEN")
	chatID := os.Getenv("CHAT_ID")

	if botToken == "" || chatID == "" {
		log.Fatal("BOT_TOKEN or CHAT_ID not set in environment")
		return fmt.Errorf("BOT_TOKEN or CHAT_ID not set")
	}

	telegramAPI := "https://api.telegram.org/bot" + botToken + "/sendMessage"

	data := url.Values{
		"chat_id":      {chatID},
		"text":         {htmlContent},
		"parse_mode":   {"HTML"},
		"disable_web_page_preview": {"true"},
	}
	resp, err := http.PostForm(telegramAPI, data)
	if err != nil {
		log.Fatal("Error sending message: ", err)
		return err
	}
	defer resp.Body.Close()

	// Check for non-OK status codes
	if resp.StatusCode != 200 {
		log.Printf("Telegram API error: Received non-OK response code %d", resp.StatusCode)
		return fmt.Errorf("telegram API returned status code: %d", resp.StatusCode)
	}

	log.Println("Report sent successfully")
	return nil
}
