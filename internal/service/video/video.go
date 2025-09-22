package video

import (
	"errors"
	"regexp"

	"github.com/bifidokk/recipe-bot/internal/service"
	"github.com/bifidokk/recipe-bot/internal/service/api"
	"github.com/rs/zerolog/log"
)

const SharedURL = "shared_url"
const VideoID = "video_id"

const TikTok = "tiktok"
const Instagram = "instagram"

type videoService struct {
	tikhub      service.TikHubClient
	instaloader service.InstaloaderClient
}

type videoIdentification struct {
	id     string
	idType string
	source string
}

func NewVideoService(
	tikhub service.TikHubClient,
	instaloader service.InstaloaderClient,
) service.VideoService {
	return &videoService{
		tikhub,
		instaloader,
	}
}

func (t *videoService) HasVideo(message string) bool {
	_, err := t.extractVideoIdentification(message)

	return err == nil
}

func (t *videoService) GetVideoData(message string) (*api.VideoData, error) {
	videoIdentificator, err := t.extractVideoIdentification(message)

	if err != nil {
		return nil, err
	}

	var videoData *api.VideoData

	switch videoIdentificator.source {
	case TikTok:
		if videoIdentificator.idType == SharedURL {
			videoData, err = t.tikhub.GetVideoDataBySharedURL(videoIdentificator.id)
		}

		if videoIdentificator.idType == VideoID {
			videoData, err = t.tikhub.GetVideoDataByVideoID(videoIdentificator.id)
		}
	case Instagram:
		videoData, err = t.instaloader.GetVideoDataBySharedURL(videoIdentificator.id)

	default:
		return nil, errors.New("unknown source")
	}

	if err != nil {
		return nil, err
	}

	videoData.Source = videoIdentificator.source
	videoData.SourceID = videoIdentificator.id
	videoData.SourceIDType = videoIdentificator.idType

	log.Info().Msgf("Video data: %v", videoData)

	return videoData, nil
}

func (t *videoService) extractVideoIdentification(message string) (*videoIdentification, error) {
	id := extractShareTikTokURL(message)

	if len(id) > 0 {
		return &videoIdentification{
			id:     id,
			idType: SharedURL,
			source: TikTok,
		}, nil
	}

	id = extractTikTokVideoID(message)

	if len(id) > 0 {
		return &videoIdentification{
			id:     id,
			idType: VideoID,
			source: TikTok,
		}, nil
	}

	id = extractInstagramVideoID(message)

	if len(id) > 0 {
		return &videoIdentification{
			id:     id,
			idType: SharedURL,
			source: Instagram,
		}, nil
	}

	return nil, errors.New("could not extract video identification")
}

func extractInstagramVideoID(message string) string {
	patterns := []string{
		`https?://(www\.)?instagram\.com/reel/([A-Za-z0-9_-]+)`,
		`https?://(www\.)?instagram\.com/p/([A-Za-z0-9_-]+)`,
		`https?://instagram\.com/reel/([A-Za-z0-9_-]+)`,
		`https?://instagram\.com/p/([A-Za-z0-9_-]+)`,
	}

	for _, pattern := range patterns {
		regex := regexp.MustCompile(pattern)
		if matches := regex.FindString(message); matches != "" {
			return matches
		}
	}

	return ""
}

func extractShareTikTokURL(message string) string {
	tikTokShareURLPattern := regexp.MustCompile(`https?://(vm|vt)\.tiktok\.com/[A-Za-z0-9]+/?`)

	return tikTokShareURLPattern.FindString(message)
}

func extractTikTokVideoID(message string) string {
	tiktokVideoIDPattern := regexp.MustCompile(`https?://(www\.)?tiktok\.com/@[A-Za-z0-9_.]+/video/([0-9]+)`)

	if matches := tiktokVideoIDPattern.FindStringSubmatch(message); len(matches) > 2 {
		return matches[2] // The second capture group is the video ID
	}

	return ""
}
