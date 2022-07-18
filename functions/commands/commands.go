package commands

import (
	"fmt"
	uni "github.com/acecca/go-unistudium"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"main/enums"
	"main/enums/callqueries"
	"main/enums/calltexts"
	"main/enums/states"
	"main/functions/cache"
	"strings"
)

var Bot *tgbotapi.BotAPI

func NewStart(message *tgbotapi.Message) {
	if _, err := cache.State(message.From.ID); err == nil {
		Start(message)
		return
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(enums.Welcome, message.From.FirstName))
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(calltexts.SearchClass, callqueries.SearchClass)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(calltexts.SearchExam, callqueries.SearchExam)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(calltexts.SearchGraduation, callqueries.SearchGraduation)),
	)

	if _, err := cache.State(message.From.ID, states.Home); err != nil {
		fmt.Println(err)
		return
	}

	if _, err := Bot.Send(msg); err != nil {
		fmt.Println(err)
		return
	}
}

func Start(message *tgbotapi.Message, fromID ...int64) {
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(enums.Start))
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(calltexts.SearchClass, callqueries.SearchClass)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(calltexts.SearchExam, callqueries.SearchExam)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(calltexts.SearchGraduation, callqueries.SearchGraduation)),
	)

	if len(fromID) == 1 {
		if _, err := cache.State(fromID[0], states.Home); err != nil {
			fmt.Println(err)
			return
		}
	} else {
		if _, err := cache.State(message.From.ID, states.Home); err != nil {
			fmt.Println(err)
			return
		}
	}

	if _, err := Bot.Send(msg); err != nil {
		fmt.Println(err)
		return
	}
}

func SearchType(query *tgbotapi.CallbackQuery, data string) {
	msg := tgbotapi.NewEditMessageTextAndMarkup(query.Message.Chat.ID, query.Message.MessageID,
		fmt.Sprintf(enums.TypeKeyword, data),
		tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(calltexts.Back, callqueries.Back))))

	state := stringToState(data)

	if _, err := cache.State(query.From.ID, state); err != nil {
		fmt.Println(err)
		return
	}

	if _, err := Bot.Send(msg); err != nil {
		fmt.Println(err)
		return
	}
}

func TypeQuery(message *tgbotapi.Message) {
	searchType, err := cache.State(message.From.ID)
	if err != nil {
		fmt.Println("Errore ottenimento stato")
		return
	}
	text := "*Risultati trovati*"
	var rooms []uni.Room
	switch searchType {
	case states.TypeClass:
		if rooms, err = uni.FindRooms(uni.ClassRoom, message.Text); err != nil {
			fmt.Println("Errore ottenimento lezioni")
			return
		}
		for i := 0; i < len(rooms); i++ {
			more := fmt.Sprintf("\n––––––––––––––––––––––––––\n\n*Nome corso*: %s\n*Prof*: %s\n*Laurea*: %s\n[Link](%s)",
				tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, rooms[i].CourseName),
				tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, rooms[i].Professor),
				tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, rooms[i].Degree),
				rooms[i].MeetingLink)
			if len(more)+len(text) < 4096 {
				text += more
				continue
			}
			msg := tgbotapi.NewMessage(message.Chat.ID, text)
			msg.DisableWebPagePreview = true
			msg.ParseMode = tgbotapi.ModeMarkdownV2
			if _, err := Bot.Send(msg); err != nil {
				fmt.Println(err)
			}
			text = ""
		}

	case states.TypeExam:
		if rooms, err = uni.FindRooms(uni.ExamRoom, message.Text); err != nil {
			fmt.Println("Errore ottenimento esami")
			return
		}
		for i := 0; i < len(rooms); i++ {
			more := fmt.Sprintf("\n––––––––––––––––––––––––––\n\n*Nome corso*: %s\n*Prof*: %s\n*Laurea*: %s\n*Data*: %s\n[Link](%s)",
				tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, rooms[i].CourseName),
				tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, rooms[i].Professor),
				tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, rooms[i].Degree),
				tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, rooms[i].Time),
				rooms[i].MeetingLink)
			if len(more)+len(text) < 4096 {
				text += more
				continue
			}
			msg := tgbotapi.NewMessage(message.Chat.ID, text)
			msg.DisableWebPagePreview = true
			msg.ParseMode = tgbotapi.ModeMarkdownV2
			if _, err := Bot.Send(msg); err != nil {
				fmt.Println(err)
			}
			text = ""
		}

	case states.TypeGraduation:
		if rooms, err = uni.FindRooms(uni.GraduationRoom, message.Text); err != nil {
			fmt.Println("Errore ottenimento lauree")
			return
		}
		for i := 0; i < len(rooms); i++ {
			more := fmt.Sprintf("\n––––––––––––––––––––––––––\n\n*Dipartimento*: %s\n*Laurea*: %s\n*Data*: %s\n*Codice*: `%s`",
				tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, rooms[i].Department),
				tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, rooms[i].Degree),
				tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, rooms[i].Time),
				rooms[i].GraduationCode)
			if len(more)+len(text) < 4096 {
				text += more
				continue
			}
			msg := tgbotapi.NewMessage(message.Chat.ID, text)
			msg.DisableWebPagePreview = true
			msg.ParseMode = tgbotapi.ModeMarkdownV2
			if _, err := Bot.Send(msg); err != nil {
				fmt.Println(err)
			}
			text = ""
		}

	}
	if strings.Contains(text, "\n") {
		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		msg.DisableWebPagePreview = true
		msg.ParseMode = tgbotapi.ModeMarkdownV2
		if _, err := Bot.Send(msg); err != nil {
			fmt.Println(err)
		}
	} else {
		_, err = Bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Nessun risultato trovato"))
		if err != nil {
			fmt.Println(err)
		}
	}

	Start(message)
}

func Back(query *tgbotapi.CallbackQuery) {
	msg := tgbotapi.NewDeleteMessage(query.Message.Chat.ID, query.Message.MessageID)
	if _, err := Bot.Request(msg); err != nil {
		fmt.Println(err)
		return
	}
	Start(query.Message, query.From.ID)
}

func stringToState(data string) uint16 {
	switch data {
	case "lezione":
		return states.TypeClass

	case "esame":
		return states.TypeExam

	case "laurea":
		return states.TypeGraduation

	default:
		return states.Home
	}
}
