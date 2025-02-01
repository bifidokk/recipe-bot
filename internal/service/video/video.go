package video

import (
	"errors"
	"regexp"

	"github.com/bifidokk/recipe-bot/internal/service"
	"github.com/bifidokk/recipe-bot/internal/service/api"
)

const SharedURL = "shared_url"
const VideoID = "video_id"

const TikTok = "tiktok"

type videoService struct {
	tikhub service.TikHubClient
}

type videoIdentification struct {
	id     string
	idType string
	source string
}

func NewVideoService(
	tikhub service.TikHubClient,
) service.VideoService {
	return &videoService{
		tikhub,
	}
}

func (t *videoService) GetVideoData(message string) (*api.VideoData, error) {
	videoIdentificator, err := t.extractVideoIdentification(message)

	if err != nil {
		return nil, err
	}

	switch videoIdentificator.source {
	case TikTok:
		if videoIdentificator.idType == SharedURL {
			return t.tikhub.GetVideoDataBySharedURL(videoIdentificator.id)
		}

		if videoIdentificator.idType == VideoID {
			return t.tikhub.GetVideoDataByVideoID(videoIdentificator.id)
		}
	}

	return nil, errors.New("invalid source")
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

	return nil, errors.New("could not extract video identification")
}

func extractShareTikTokURL(s string) string {
	tikTokShareURLPattern := regexp.MustCompile(`https?://(vm|vt)\.tiktok\.com/[A-Za-z0-9]+/?`)

	return tikTokShareURLPattern.FindString(s)
}

func extractTikTokVideoID(message string) string {
	tiktokVideoIDPattern := regexp.MustCompile(`https?://(www\.)?tiktok\.com/@[A-Za-z0-9_.]+/video/([0-9]+)`)

	if matches := tiktokVideoIDPattern.FindStringSubmatch(message); len(matches) > 2 {
		return matches[2] // The second capture group is the video ID
	}

	return ""
}
