package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	url := "http://172.28.0.1:8090/"

	targetTime := time.Date(2024, time.Now().Month(), time.Now().Day(), 18, 16, 0, 0, time.Local)

	for {
		now := time.Now()

		if now.After(targetTime) {
			resp, err := http.Get(url)
			if err != nil {
				fmt.Printf("Error: ", err)
				return
			}
			defer resp.Body.Close()

			//TODO:
			// parse the html...
			// register the laundry automatically (press the button)

			break
		}

		time.Sleep(10 * time.Second)
	}
}
