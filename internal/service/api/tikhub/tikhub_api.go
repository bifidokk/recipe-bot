package tikhub

import (
	"encoding/json"
	"fmt"
	"github.com/bifidokk/recipe-bot/internal/service/api"
	"io"
	"net/http"
)

const tikTokGetVideoByShareURL = "https://beta.tikhub.io/api/v1/tiktok/app/v3/fetch_one_video_by_share_url?share_url="
const tikTokGetVideoById = "https://beta.tikhub.io/api/v1/tiktok/app/v3/fetch_one_video?aweme_id="

type VideoDataResponse struct {
	Data struct {
		AwemeDetails []struct {
			Music struct {
				PlayURL struct {
					URI string `json:"uri"`
				} `json:"play_url"`
			} `json:"music"`
			OriginalClientText struct {
				MarkupText string `json:"markup_text"`
			} `json:"original_client_text"`
		} `json:"aweme_details"`
	} `json:"data"`
}

type TikHubClient struct {
	httpClient *http.Client
	apiToken   string
}

func NewTikHubClient(apiToken string) *TikHubClient {
	client := &http.Client{}

	return &TikHubClient{
		httpClient: client,
		apiToken:   apiToken,
	}
}

func (t *TikHubClient) GetVideoDataBySharedURL(sharedURL string) (*api.VideoData, error) {
	response, err := t.request("GET", tikTokGetVideoByShareURL+sharedURL, nil)

	if err != nil {
		return nil, err
	}

	videoDataResponse := &VideoDataResponse{}

	err = json.Unmarshal(response, videoDataResponse)
	if err != nil {
		return nil, err
	}

	return &api.VideoData{
		AudioURL:    videoDataResponse.Data.AwemeDetails[0].Music.PlayURL.URI,
		Description: videoDataResponse.Data.AwemeDetails[0].OriginalClientText.MarkupText,
	}, nil
}

func (t *TikHubClient) GetVideoDataByVideoId(videoId string) (*api.VideoData, error) {
	response, err := t.request("GET", tikTokGetVideoById+videoId, nil)

	if err != nil {
		return nil, err
	}

	videoDataResponse := &VideoDataResponse{}

	err = json.Unmarshal(response, videoDataResponse)
	if err != nil {
		return nil, err
	}

	return &api.VideoData{
		AudioURL:    videoDataResponse.Data.AwemeDetails[0].Music.PlayURL.URI,
		Description: videoDataResponse.Data.AwemeDetails[0].OriginalClientText.MarkupText,
	}, nil

}

func (t *TikHubClient) request(method string, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+t.apiToken)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data: %s", resp.Status)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, err
}
