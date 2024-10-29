package bot

import (
	"github.com/bifidokk/recipe-bot/internal/service"
	"github.com/bifidokk/recipe-bot/internal/service/utils"
	"gopkg.in/tucnak/telebot.v2"
	"log"
)

type botService struct {
	bot          *telebot.Bot
	apiToken     string
	openai       service.OpenAIClient
	videoService service.VideoService
}

func NewBotService(
	bot *telebot.Bot,
	apiToken string,
	openai service.OpenAIClient,
	videoService service.VideoService,
) service.BotService {
	return &botService{
		bot,
		apiToken,
		openai,
		videoService,
	}
}

func (bs *botService) Start() error {
	log.Println("Starting bot")

	bs.bot.Handle(telebot.OnText, bs.onTextMessage)

	log.Println("Bot started!")
	go func() {
		bs.bot.Start()
	}()

	return nil
}

func (bs *botService) onTextMessage(message *telebot.Message) {
	log.Println(message.Text)

	videoData, err := bs.videoService.GetVideoData(message.Text)

	if err != nil {
		log.Println(err)
	}

	log.Println(videoData)

	filePath, err := utils.DownloadFileFromURL(videoData.AudioURL)

	if err != nil {
		log.Println(err)
		return
	}

	text, err := bs.openai.ConvertSpeechToText(filePath)

	if err != nil {
		log.Println(err)
		return
	}

	recipeText, err := bs.openai.TextToFormattedRecipe(text, videoData.Description)

	bs.bot.Send(message.Sender, recipeText)
}
