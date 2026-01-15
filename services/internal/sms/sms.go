package sms

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
)

func SendSMS(to, message string) error {
	username := os.Getenv("AT_USERNAME")
	url := os.Getenv("AT_SMS_URL")
	key := os.Getenv("AT_API_KEY")
	shortcode := os.Getenv("AT_SHORTCODE")

	body := []byte(fmt.Sprintf("username=%s&to=%s&message=%s&from=%s", username, to, message, shortcode))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("API Key %s", key))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Handle the response here
	// fmt.Println(resp.StatusCode)
	log.Println(resp.Status)

	return nil
}
