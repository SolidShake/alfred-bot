package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/SolidShake/alfred-bot/internal/api/bank"
	"github.com/SolidShake/alfred-bot/internal/api/binance"
	"github.com/SolidShake/alfred-bot/internal/api/hackernews"
	b "github.com/SolidShake/alfred-bot/internal/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

//@TODO add linter
func main() {
	if os.Getenv("APP_ENV") != "prod" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("error loading .env file")
		}
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TGBOT_API_KEY"))
	if err != nil {
		fmt.Println(err)
		return
	}

	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err == nil {
		bot.Debug = debug
	}

	fullResponse := b.FullResponse{
		Blocks: []string{
			"***Доброе утро*** \xF0\x9F\x8C\x9E",
			bank.CreateBankResponse(),
			binance.GetResponse(),
			hackernews.GetResponse(),
		},
	}

	chatID, _ := strconv.ParseInt(os.Getenv("DEBUG_CHAT_ID"), 10, 64)

	ticker := time.NewTicker(5 * time.Second)
	stop := make(chan struct{})
	defer close(stop)
	for {
		select {
		case <-ticker.C:
			msg := tgbotapi.NewMessage(chatID, fullResponse.ToString())
			msg.ParseMode = "markdown"
			msg.DisableWebPagePreview = true
			bot.Send(msg)
		case <-stop:
			ticker.Stop()
			return
		}
	}
}
