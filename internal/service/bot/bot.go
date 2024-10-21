package bot

import (
	"github.com/bifidokk/receipt-bot/internal/service"
	"gopkg.in/tucnak/telebot.v2"
	"log"
)

type botService struct {
	bot *telebot.Bot
}

func NewBotService(bot *telebot.Bot) service.BotService {
	return &botService{bot}
}

func (bs *botService) Start() error {
	log.Println("Starting bot")

	bs.bot.Handle(telebot.OnText, resolveMessage)

	log.Println("Bot started!")
	go func() {
		bs.bot.Start()
	}()

	return nil
}

func resolveMessage(message *telebot.Message) {
	log.Println(message.Text)
}
