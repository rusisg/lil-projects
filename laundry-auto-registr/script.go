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

	urlToCheck := os.Getenv("URL")                  // The URL to check
	buttonPressURL := os.Getenv("BUTTON_PRESS_URL") // The form submission endpoint

	timeZone := time.FixedZone("GMT+1", 1*60*60)
	targetTime := time.Date(2024, time.Now().Month(), time.Now().Day(), 17, 16, 0, 0, timeZone)
	//targetTime := time.Date(2024, time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), 0, 0, timeZone)

	for {
		now := time.Now()

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
