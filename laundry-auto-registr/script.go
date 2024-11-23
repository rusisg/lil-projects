package main

import (
	"bytes"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	url := os.Getenv("URL")                         // The URL to check
	buttonPressURL := os.Getenv("BUTTON_PRESS_URL") // The API endpoint for pressing the button

	timeZone := time.FixedZone("GMT+1", 1*60*60)
	targetTime := time.Date(2024, time.Now().Month(), time.Now().Day(), 18, 16, 0, 0, timeZone)

	for {
		now := time.Now()
		log.Println("Starting the loop...")

		if now.After(targetTime) {
			// Check the status of the URL
			resp, err := http.Get(url)
			fmt.Println("Get URL")
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("Body close")
			}(resp.Body)

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Error while reading -> %v\n", err)
				return
			}
			fmt.Printf("Received: %s\n", body)

			// Register the laundry automatically (press the button)
			requestBody := []byte(`{"action": "press"}`) // Replace with the actual request body if needed

			req, err := http.NewRequest("POST", buttonPressURL, bytes.NewBuffer(requestBody))
			if err != nil {
				fmt.Printf("Error creating request: %v\n", err)
				return
			}

			client := &http.Client{}
			resp, err = client.Do(req)
			if err != nil {
				fmt.Printf("Error sending request: %v\n", err)
				return
			}

			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					log.Fatal(err)
				}
			}(resp.Body)

			responseBody, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Error reading response: %v\n", err)
				return
			}
			fmt.Printf("Button press response: %s\n", responseBody)

			break
		}

		time.Sleep(10 * time.Second)
	}
}
