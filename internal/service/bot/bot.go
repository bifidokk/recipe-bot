package bot

import (
	"github.com/bifidokk/recipe-bot/internal/entity"
	"github.com/bifidokk/recipe-bot/internal/service"
	"github.com/bifidokk/recipe-bot/internal/service/recipe"
	"github.com/bifidokk/recipe-bot/internal/service/utils"
	"github.com/rs/zerolog/log"
	"gopkg.in/telebot.v4"
)

type botService struct {
	bot           *telebot.Bot
	apiToken      string
	openai        service.OpenAIClient
	videoService  service.VideoService
	recipeService recipe.Service
}

func NewBotService(
	bot *telebot.Bot,
	apiToken string,
	openai service.OpenAIClient,
	videoService service.VideoService,
	recipeService recipe.Service,
) service.BotService {
	return &botService{
		bot,
		apiToken,
		openai,
		videoService,
		recipeService,
	}
}

func (bs *botService) Start() error {
	log.Info().Msg("Starting bot")

	bs.bot.Handle(telebot.OnText, bs.onTextMessage)

	log.Info().Msg("Bot started!")
	go func() {
		bs.bot.Start()
	}()

	return nil
}

func (bs *botService) onTextMessage(c telebot.Context) error {
	log.Info().Msgf("Input text %v", c.Text())

	videoData, err := bs.videoService.GetVideoData(c.Text())

	if err != nil {
		_, err = bs.bot.Send(c.Sender(), "Sorry but I could not get video data from your message")

		log.Error().Err(err)
		return err
	}

	log.Info().Msgf("Video data: %v", videoData)

	recipeData := &recipe.CreateRecipeData{
		AudioURL:     videoData.AudioURL,
		Source:       videoData.Source,
		SourceID:     videoData.SourceID,
		SourceIDType: videoData.SourceIDType,
		ShareURL:     videoData.ShareURL,
	}

	filePath, err := utils.DownloadFileFromURL(videoData.AudioURL)

	if err != nil {
		log.Error().Err(err).Msg("Failed to download video file")
		return err
	}

	text, err := bs.openai.ConvertSpeechToText(filePath)

	if err != nil {
		log.Error().Err(err).Msg("Failed to convert speech to text")
		return err
	}

	recipeTextData, err := bs.openai.TextToFormattedRecipe(text, videoData.Description)

	if err != nil {
		log.Error().Err(err).Msg("Failed to convert text to recipe")
		return err
	}

	recipeData.Body = text
	recipeData.RecipeMarkdownText = recipeTextData.Text
	recipeData.Title = recipeTextData.Title

	u := c.Get("user").(*entity.User)

	r, err := bs.recipeService.CreateRecipe(
		recipeData,
		u.ID,
	)

	if err != nil {
		log.Error().Err(err).Msg("Failed to create recipe")
		return err
	}

	if videoData.CoverURL != "" {
		photo := &telebot.Photo{File: telebot.FromURL(videoData.CoverURL)}
		photoResult, _ := bs.bot.Send(c.Sender(), photo)

		if photoResult != nil {
			r.CoverFileID = photoResult.Photo.FileID
			err = bs.recipeService.UpdateRecipe(r)

			if err != nil {
				log.Error().Err(err).Msg("Failed to update recipe")
			}
		}
	}

	_, err = bs.bot.Send(c.Sender(), r.RecipeMarkdownText)

	if err != nil {
		log.Error().Err(err).Msg("Failed to send message")
		return err
	}

	return nil
}
