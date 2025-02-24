package service

import (
	"github.com/bifidokk/recipe-bot/internal/entity"
	"github.com/bifidokk/recipe-bot/internal/service/api"
	"github.com/bifidokk/recipe-bot/internal/service/api/openai"
)

type BotService interface {
	Start() error
}

type OpenAIClient interface {
	ConvertSpeechToText(inputFile string) (string, error)
	TextToFormattedRecipe(speechText string, descriptionText string) (*openai.Recipe, error)
}

type TikHubClient interface {
	GetVideoDataBySharedURL(sharedURL string) (*api.VideoData, error)
	GetVideoDataByVideoID(videoID string) (*api.VideoData, error)
}

type VideoService interface {
	GetVideoData(message string) (*api.VideoData, error)
}

type UserService interface {
	GetUserByTelegramID(ID int64) (*entity.User, error)
	CreateUser(user *entity.User) (*entity.User, error)
}
