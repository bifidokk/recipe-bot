package openai

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
)

const speechToTextPromt string = "Use capital letters and punctuation. Do not repeat yourself. Do not describe ambient sounds or noise or silence, just ommit. Break the text into paragraphs. Separate paragraphs with blank lines"

const textToRecipePromt = `
You are a helpful assistant in cooking
In input, you receive two messages.
The Message1 is speech converted to text that describes food recipes. It usually contains an explanation of the cooking process.
The Message2 is an additional description of this recipe and contains details such as ingredients and their amount, time of cooking etc.
Create a recipe title and text in the JSON format. Return raw JSON without formatting. The format:
###
{
	"title": "Recipe title",
	"text": "Recipe text",
}
###

The first element with key "title" - just create a small title for the recipe, not more that 255 symbols.
The second element with key "text" - recipe text in format:
###
Name of dish.
List of ingredients.
Cooking process description.
###

The List of ingredients in recipe text should be formatted as list. Blocks in recipe text should be divided by blank line.
The cooking process description should be formatted and divided by steps.

Answer using a language from input messages
Message1:
###
{message1}
###

Message2:
###
{message2}
###
`

type Client struct {
	client *openai.Client
}

type Recipe struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

func NewOpenAIClient(token string) *Client {
	client := openai.NewClient(token)

	return &Client{
		client: client,
	}
}

func (c Client) ConvertSpeechToText(inputFile string) (string, error) {
	ctx := context.Background()

	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: inputFile,
		Prompt:   speechToTextPromt,
	}

	resp, err := c.client.CreateTranscription(ctx, req)

	if err != nil {
		log.Error().Err(err).Msg("Transcription error")
		return "", err
	}

	log.Info().Msgf("Result text: %v", resp.Text)

	return resp.Text, nil
}

func (c Client) TextToFormattedRecipe(speechText string, descriptionText string) (*Recipe, error) {
	ctx := context.Background()

	message := strings.Replace(textToRecipePromt, "{message1}", speechText, 1)
	message = strings.Replace(message, "{message2}", descriptionText, 1)

	req := openai.ChatCompletionRequest{
		Model: openai.GPT4o,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: message,
			},
		},
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)

	if err != nil {
		log.Error().Err(err).Msg("Open AI chat completion error")
		return nil, err
	}

	log.Info().Msgf("Chat completion result: %v", resp.Choices[0].Message.Content)

	jsonData := []byte(resp.Choices[0].Message.Content)

	recipe, err := parseRecipeResponse(jsonData)
	if err != nil {
		return nil, err
	}

	return recipe, nil
}

func parseRecipeResponse(data []byte) (*Recipe, error) {
	var recipe Recipe
	err := json.Unmarshal(data, &recipe)
	if err != nil {
		return nil, err
	}

	return &recipe, nil
}
