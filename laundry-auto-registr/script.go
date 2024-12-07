package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	urlToCheck := os.Getenv("URL")                  // The URL to check
	buttonPressURL := os.Getenv("BUTTON_PRESS_URL") // The form submission endpoint

	timeZone := time.FixedZone("GMT+1", 1*60*60)
	targetTime := time.Date(2024, time.Now().Month(), time.Now().Day(), 17, 16, 0, 0, timeZone)
	// targetTime := time.Date(2024, time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), 0, 0, timeZone)

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
	defer cancel()

	for {
		now := time.Now()

		select {
		case <-ctx.Done():
			fmt.Println("Context deadline exceeded. Exiting loop...")
			return
		default:
		}
		if now.After(targetTime) {
			// Check the status of the URL
			resp, err := http.Get(urlToCheck)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			defer resp.Body.Close()

			// Submit the form to press the button
			formData := url.Values{}
			formData.Set("reg", "20241125-08:00-8") // Replace with the actual value from your HTML

			resp, err = http.PostForm(buttonPressURL, formData)
			if err != nil {
				fmt.Printf("Error sending request: %v\n", err)
				return
			}
			defer resp.Body.Close()

			fmt.Printf("Button press response status: %s\n", resp.Status)

			break
		}

		time.Sleep(10 * time.Second)
	}
}
