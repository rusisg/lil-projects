package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	urlToCheck := os.Getenv("URL")
	buttonPressURL := os.Getenv("BUTTON_PRESS_URL")

	// Define the target time range
	timeZone := time.FixedZone("GMT+1", 1*60*60)
	timeStart := time.Date(2024, time.Now().Month(), time.Now().Day(), 17, 47, 0, 0, timeZone)
	timeStop := time.Date(2024, time.Now().Month(), time.Now().Day(), 17, 50, 0, 0, timeZone)

	// Wait until timeStart
	for time.Now().Before(timeStart) {
		time.Sleep(1 * time.Second)
		fmt.Printf("Waiting for the start time: %s...\n", timeStart)
	}

	// Execute the loop between timeStart and timeStop
	for time.Now().Before(timeStop) {
		// Perform the GET request
		resp, err := http.Get(urlToCheck)
		if err != nil {
			fmt.Printf("Error performing GET request: %v\n", err)
			return
		}
		resp.Body.Close() // Close response body to avoid resource leaks

		// Perform the POST request
		formData := url.Values{}
		formData.Set("reg", "20241125-08:00-8")

		resp, err = http.PostForm(buttonPressURL, formData)
		if err != nil {
			fmt.Printf("Error sending POST request: %v\n", err)
			return
		}
		defer resp.Body.Close()

		fmt.Printf("Button press response status: %s\n", resp.Status)

		// Exit after successfully pressing the button
		break
	}

	// If the script reaches here, the operation is complete or timeStop has been exceeded
	fmt.Println("Operation completed or time range expired.")
}
