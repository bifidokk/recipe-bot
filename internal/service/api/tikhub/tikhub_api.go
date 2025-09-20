package tikhub

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/bifidokk/recipe-bot/internal/service/api"
)

const tikTokGetVideoByShareURL = "https://beta.tikhub.io/api/v1/tiktok/app/v3/fetch_one_video_by_share_url?share_url="
const tikTokGetVideoByID = "https://beta.tikhub.io/api/v1/tiktok/app/v3/fetch_one_video?aweme_id="

type VideoDataResponse struct {
	Code int `json:"code"`
	Data struct {
		AwemeDetail struct {
			Music struct {
				PlayURL struct {
					URI     string   `json:"uri"`
					URLList []string `json:"url_list"`
				} `json:"play_url"`
			} `json:"music"`
			OriginalClientText struct {
				MarkupText string `json:"markup_text"`
			} `json:"original_client_text"`
			Video struct {
				AIDynamicCover struct {
					URI     string   `json:"uri"`
					URLList []string `json:"url_list"`
				} `json:"ai_dynamic_cover"`
			} `json:"video"`
			ShareURL string `json:"share_url"`
		} `json:"aweme_detail"`
	} `json:"data"`
}

type Client struct {
	httpClient *http.Client
	apiToken   string
}

func NewTikHubClient(apiToken string) *Client {
	client := &http.Client{}

	return &Client{
		httpClient: client,
		apiToken:   apiToken,
	}
}

func (t *Client) GetVideoDataBySharedURL(sharedURL string) (*api.VideoData, error) {
	response, err := t.request("GET", tikTokGetVideoByShareURL+sharedURL, nil)

	if err != nil {
		return nil, err
	}

	log.Info().Msgf("TikHub API Response (ShareURL): %s", string(response))

	videoDataResponse := &VideoDataResponse{}

	err = json.Unmarshal(response, videoDataResponse)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to unmarshal TikHub response: %s", string(response))
		return nil, err
	}

	if videoDataResponse.Code != 200 {
		log.Error().Msgf("TikHub API returned non-success code: %d", videoDataResponse.Code)
		return nil, fmt.Errorf("TikHub API returned error code: %d", videoDataResponse.Code)
	}

	videoData := &api.VideoData{
		AudioURL:    videoDataResponse.Data.AwemeDetail.Music.PlayURL.URI,
		Description: videoDataResponse.Data.AwemeDetail.OriginalClientText.MarkupText,
		ShareURL:    videoDataResponse.Data.AwemeDetail.ShareURL,
	}

	// Safely get cover URL
	if len(videoDataResponse.Data.AwemeDetail.Video.AIDynamicCover.URLList) > 0 {
		videoData.CoverURL = videoDataResponse.Data.AwemeDetail.Video.AIDynamicCover.URLList[0]
	}

	return videoData, nil
}

func (t *Client) GetVideoDataByVideoID(videoID string) (*api.VideoData, error) {
	response, err := t.request("GET", tikTokGetVideoByID+videoID, nil)

	if err != nil {
		return nil, err
	}

	log.Info().Msgf("TikHub API Response (VideoID): %s", string(response))

	videoDataResponse := &VideoDataResponse{}

	err = json.Unmarshal(response, videoDataResponse)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to unmarshal TikHub response: %s", string(response))
		return nil, err
	}

	if videoDataResponse.Code != 200 {
		log.Error().Msgf("TikHub API returned non-success code: %d", videoDataResponse.Code)
		return nil, fmt.Errorf("TikHub API returned error code: %d", videoDataResponse.Code)
	}

	log.Info().Msgf("Video data: %v", videoDataResponse)

	videoData := &api.VideoData{
		AudioURL:    videoDataResponse.Data.AwemeDetail.Music.PlayURL.URI,
		Description: videoDataResponse.Data.AwemeDetail.OriginalClientText.MarkupText,
		ShareURL:    videoDataResponse.Data.AwemeDetail.ShareURL,
	}

	// Safely get cover URL
	if len(videoDataResponse.Data.AwemeDetail.Video.AIDynamicCover.URLList) > 0 {
		videoData.CoverURL = videoDataResponse.Data.AwemeDetail.Video.AIDynamicCover.URLList[0]
	}

	return videoData, nil

}

func (t *Client) request(method string, url string, body io.Reader) ([]byte, error) {
	log.Info().Msgf("Making TikHub API request to: %s", url)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+t.apiToken)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("TikHub API request failed")
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	log.Info().Msgf("TikHub API response status: %s", resp.Status)

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		log.Error().Msgf("TikHub API error response: %s", string(responseBody))
		return nil, fmt.Errorf("failed to fetch data: %s", resp.Status)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read TikHub API response body")
		return nil, err
	}

	return responseBody, err
}
