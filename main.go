package main

import (
	_ "embed"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"main/enums/callqueries"
	"main/enums/states"
	"main/functions/cache"
	"main/functions/commands"
	"os"
	"strings"
)

//go:embed token
var token string

func main() {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	commands.Bot = bot

	up := tgbotapi.NewUpdate(0)
	up.AllowedUpdates = []string{"message", "callback_query"}
	updates := bot.GetUpdatesChan(up)

	for up := range updates {
		go func(update tgbotapi.Update) {
			defer func() {
				if err := recover(); err != nil {
					log.Println(err)
				}
			}()
			if update.Message != nil {
				msg := update.Message
				if msg.IsCommand() {
					if msg.Command() == "start" {
						commands.NewStart(msg)
						return
					}
				}
				state, err := cache.State(msg.From.ID)
				if err != nil {
					return
				}
				if state == states.TypeClass || state == states.TypeExam || state == states.TypeGraduation {
					commands.TypeQuery(msg)
				}
			} else if update.CallbackQuery != nil {
				query := update.CallbackQuery
				data := strings.Split(query.Data, "#")
				switch data[0] {
				case callqueries.SearchType:
					commands.SearchType(query, data[1])

				case callqueries.Back:
					commands.Back(query)
				}
			}
		}(up)
	}
}
