package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	url := os.Getenv("URL")

	targetTime := time.Date(2024, time.Now().Month(), time.Now().Day(), 18, 16, 0, 0, time.UTC)

	for {
		now := time.Now()

		if now.After(targetTime) {
			resp, err := http.Get(url)
			if err != nil {
				fmt.Printf("Error: ", err)
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Error while reading -> %v\n", err)
				return
			}
			fmt.Printf("Received %s\n", body)
			//TODO:
			// parse the html...
			// register the laundry automatically (press the button)

			break
		}

		time.Sleep(10 * time.Second)
	}
}
