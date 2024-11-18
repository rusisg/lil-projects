package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI
var chatID int64

// Placeholder for chosen payment gateway API endpoint
var paymentGatewayURL string

// Placeholder for chosen payment gateway client ID/secret
var paymentGatewayClientID, paymentGatewayClientSecret string

func sendMessage(msg string) {
	msgConfig := tgbotapi.NewMessage(chatID, msg)
	bot.Send(msgConfig)
}

func createPaymentLink(orderDetails string) (string, error) {
	// Implement logic to create a payment link using your chosen payment gateway API
	// This may involve making an API call with order details, user data, etc.
	// Replace with actual API call logic specific to your chosen gateway

	// Placeholder example - replace with appropriate URL and request body construction

	url := fmt.Sprintf(paymentGatewayURL, paymentGatewayClientID, orderDetails)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var paymentLinkResponse map[string]interface{}
	err = json.Unmarshal(body, &paymentLinkResponse)
	if err != nil {
		return "", err
	}

	if paymentLink, ok := paymentLinkResponse["payment_link"]; ok {
		return paymentLink.(string), nil
	}
	return "", fmt.Errorf("payment link not found in response: %s", string(body))
}

func handleStart(update tgbotapi.Update) {
	if update.Message == nil || update.Message.Text != "/start" {
		return
	}

	chatID = update.Message.Chat.ID
	sendMessage("Welcome to Pizzastan! Use /order to start your pizza journey.")
}

func gatherOrderDetails(update tgbotapi.Update) string {
	// Create an inline keyboard with pizza options, sizes, and toppings
	var keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Margherita", "margherita"),
			tgbotapi.NewInlineKeyboardButtonData("Pepperoni", "pepperoni"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Small", "small"),
			tgbotapi.NewInlineKeyboardButtonData("Medium", "medium"),
			tgbotapi.NewInlineKeyboardButtonData("Large", "large"),
		),
		// ... add more rows for toppings, etc.
	)

	// Send the keyboard to the user
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please select your pizza:")
	msg.ReplyMarkup = keyboard
	bot.Send(msg)

	// Wait for user's response (you'll need to implement a loop to listen for updates)
	// Process the user's selections and build the order details string
	// ...

	// For simplicity, let's assume the user's selections are stored in variables:
	pizzaType := "margherita"
	pizzaSize := "medium"
	// ... other selections

	orderDetails := fmt.Sprintf("Pizza: %s, Size: %s, ...", pizzaType, pizzaSize)
	return orderDetails
}

func handleOrder(update tgbotapi.Update) {
	if update.Message == nil || update.Message.Text != "/order" {
		return
	}

	chatID = update.Message.Chat.ID

	// Gather order details (pizza selection, size, toppings, etc.)
	// Use Telegram Inline Keyboard or other interactive elements to capture user choices.
	// Example using Inline Keyboard (replace with your specific logic):
	orderDetails := gatherOrderDetails(update)

	paymentLink, err := createPaymentLink(orderDetails)
	if err != nil {
		sendMessage(fmt.Sprintf("Error creating payment link: %v", err))
		return
	}

	sendMessage(fmt.Sprintf("Here's your order confirmation:\n%s\nPlease click the link to complete your payment with FreedomBank:\n%s", orderDetails, paymentLink))
}

func main() {
	bot, err := tgbotapi.NewBotAPI("BOT_TOKEN")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		handleStart(update) // Assuming functionality for "/start" command
		handleOrder(update)
	}
}
