package video

import (
	"errors"
	"github.com/bifidokk/recipe-bot/internal/service"
	"github.com/bifidokk/recipe-bot/internal/service/api"
	"regexp"
)

const SharedURL = "shared_url"
const VideoId = "video_id"

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

		if videoIdentificator.idType == VideoId {
			return t.tikhub.GetVideoDataByVideoId(videoIdentificator.id)
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

	id = extractTikTokVideoId(message)

	if len(id) > 0 {
		return &videoIdentification{
			id:     id,
			idType: VideoId,
			source: TikTok,
		}, nil
	}

	return nil, errors.New("could not extract video identification")
}

func extractShareTikTokURL(s string) string {
	tiktokShareURLPattern := `https?://(vm|vt)\.tiktok\.com/[A-Za-z0-9]+/?`
	re := regexp.MustCompile(tiktokShareURLPattern)
	match := re.FindString(s)

	return match
}

func extractTikTokVideoId(message string) string {
	tiktokVideoIDPattern := `https?://(www\.)?tiktok\.com/@[A-Za-z0-9_.]+/video/([0-9]+)`
	re := regexp.MustCompile(tiktokVideoIDPattern)
	matches := re.FindStringSubmatch(message)

	if len(matches) > 2 {
		return matches[2] // The second capture group is the video ID
	}

	return ""
}
