package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/buger/jsonparser"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	telegramToken = "5832278739:AAEYuWcpxVkro1DhChmuvEe_RsTCum7N6xg"
	openWeatherAPIKey = "c5dd51445cb396d147cef472342040a9"
)

func fetchWeatherForecast() (string, error) {
	endpointURL := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=Uzhgorod,ua&units=metric&APPID=%s", openWeatherAPIKey)
	response, err := http.Get(endpointURL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	weatherDescription, _ := jsonparser.GetString(data, "weather", "[0]", "description")
	temp, _ := jsonparser.GetFloat(data, "main", "temp")

	return fmt.Sprintf("Weather: %s\nTemperature: %.2f°C", weatherDescription, temp), nil
}

func main() {
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updates, _ := bot.GetUpdatesChan(tgbotapi.UpdateConfig{Timeout: 1800})

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		forecast, err := fetchWeatherForecast()
		if err != nil {
			log.Printf("Error fetching weather: %v", err)
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, forecast)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}