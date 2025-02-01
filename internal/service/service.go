package service

import (
	"github.com/bifidokk/recipe-bot/internal/service/api"
)

type BotService interface {
	Start() error
}

type OpenAIClient interface {
	ConvertSpeechToText(inputFile string) (string, error)
	TextToFormattedRecipe(speechText string, descriptionText string) (string, error)
}

type TikHubClient interface {
	GetVideoDataBySharedURL(sharedURL string) (*api.VideoData, error)
	GetVideoDataByVideoID(videoID string) (*api.VideoData, error)
}

type VideoService interface {
	GetVideoData(message string) (*api.VideoData, error)
}
