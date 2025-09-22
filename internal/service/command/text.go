package command

import (
	"database/sql"
	"os"
	"strings"

	"github.com/bifidokk/recipe-bot/internal/entity"
	"github.com/bifidokk/recipe-bot/internal/service"
	"github.com/bifidokk/recipe-bot/internal/service/recipe"
	"github.com/bifidokk/recipe-bot/internal/service/utils"
	"github.com/rs/zerolog/log"
	"gopkg.in/telebot.v4"
)

type TextCommand struct {
	openai        service.OpenAIClient
	videoService  service.VideoService
	recipeService recipe.Service
	userService   service.UserService
}

func NewTextCommand(
	openai service.OpenAIClient,
	videoService service.VideoService,
	recipeService recipe.Service,
	userService service.UserService,
) Command {
	return &TextCommand{
		openai:        openai,
		videoService:  videoService,
		recipeService: recipeService,
		userService:   userService,
	}
}

func (c *TextCommand) Name() string {
	return "text"
}

func (c *TextCommand) Register(b *telebot.Bot) {
	b.Handle(telebot.OnText, func(ctx telebot.Context) error {
		log.Info().Msgf("Input text %v", ctx.Text())

		if !c.videoService.HasVideo(ctx.Text()) {
			_, err := b.Send(ctx.Sender(), "Sorry but I could not find video in your message")

			log.Error().Err(err)
			return err
		}

		u := ctx.Get("user").(*entity.User)

		if u.RecipeLimit <= 0 {
			_, err := b.Send(ctx.Sender(), "Sorry, you have run out of available recipes. Please check your recipe limit.")

			log.Error().Err(err)
			return err
		}

		err := c.userService.DecreaseUserLimit(u)
		if err != nil {
			return err
		}

		videoData, err := c.videoService.GetVideoData(ctx.Text())

		if err != nil {
			_, err = b.Send(ctx.Sender(), "Sorry but I could not get video data from your message")

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

		var filePath string

		// Check if AudioURL is a local file path or remote URL
		if strings.HasPrefix(videoData.AudioURL, "http://") || strings.HasPrefix(videoData.AudioURL, "https://") {
			filePath, err = utils.DownloadFileFromURL(videoData.AudioURL)
			if err != nil {
				log.Error().Err(err).Msg("Failed to download video file")
				return err
			}
		} else {
			filePath = videoData.AudioURL
			log.Info().Msgf("Using local audio file: %s", filePath)
		}

		text, err := c.openai.ConvertSpeechToText(filePath)

		if err != nil {
			log.Error().Err(err).Msg("Failed to convert speech to text")
			return err
		}

		recipeTextData, err := c.openai.TextToFormattedRecipe(text, videoData.Description)

		if err != nil {
			log.Error().Err(err).Msg("Failed to convert text to recipe")
			return err
		}

		recipeData.Body = text
		recipeData.RecipeMarkdownText = recipeTextData.Text
		recipeData.Title = recipeTextData.Title

		r, err := c.recipeService.CreateRecipe(
			recipeData,
			u.ID,
		)

		if err != nil {
			log.Error().Err(err).Msg("Failed to create recipe")
			return err
		}

		if videoData.CoverURL != "" {
			photo := &telebot.Photo{File: telebot.FromURL(videoData.CoverURL)}
			photoResult, _ := b.Send(ctx.Sender(), photo)

			if photoResult != nil {
				r.CoverFileID = sql.NullString{String: photoResult.Photo.FileID, Valid: true}
				err = c.recipeService.UpdateRecipe(r)

				if err != nil {
					log.Error().Err(err).Msg("Failed to update recipe")
				}
			}
		}

		_, err = b.Send(ctx.Sender(), r.RecipeMarkdownText)

		if err != nil {
			log.Error().Err(err).Msg("Failed to send message")
			return err
		}

		// Clean up local audio files after processing
		if !strings.HasPrefix(filePath, "http") {
			if err := os.Remove(filePath); err != nil {
				log.Warn().Err(err).Msgf("Failed to clean up local audio file: %s", filePath)
			} else {
				log.Info().Msgf("Cleaned up local audio file: %s", filePath)
			}
		}

		return nil
	})
}
