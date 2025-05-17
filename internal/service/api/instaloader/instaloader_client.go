package instaloader

import (
	"errors"
	"os/exec"
	"regexp"

	"github.com/bifidokk/recipe-bot/internal/service/api"
	"github.com/rs/zerolog/log"
)

type Client struct {
}

func NewInstaloaderClient() *Client {
	return &Client{}
}

func (t *Client) GetVideoDataBySharedURL(sharedURL string) (*api.VideoData, error) {
	pattern := regexp.MustCompile(`https?://www\.instagram\.com/reel/([a-zA-Z0-9_-]+)`)
	matches := pattern.FindStringSubmatch(sharedURL)

	if len(matches) < 2 {
		return nil, errors.New("invalid Instagram shared URL")
	}

	videoID := matches[1]

	log.Info().Msgf("Extracted video ID: %s", videoID)

	// #nosec G204 -- shortcode is validated via regex
	cmd := exec.Command("instaloader", "--dirname-pattern=/tmp/"+videoID, "--", "-"+videoID)

	if err := cmd.Run(); err != nil {
		log.Error().Err(err).Msg("Failed to download reel")
	}

	log.Info().Msgf("Downloaded video with ID: %s", videoID)

	return nil, errors.New("not implemented")
}
