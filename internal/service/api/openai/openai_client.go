package openai

import (
	"context"
	"log"
	"strings"

	"github.com/sashabaranov/go-openai"
)

const speechToTextPromt string = "Use capital letters and punctuation. Do not repeat yourself. Do not describe ambient sounds or noise or silence, just ommit. Break the text into paragraphs. Separate paragraphs with blank lines"

const textToRecipePromt = `
You are a helpful assistant in cooking
In input, you receive two messages.
The Message1 is speech converted to text that describes food recipes. It usually contains an explanation of the cooking process.
The Message2 is an additional description of this recipe and contains details such as ingredients and their amount, time of cooking etc.
Create a recipe text in the format:
###
Name of dish.
List of ingredients.
Cooking process description.
###

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
		log.Printf("Transcription error: %v\n", err)
		return "", err
	}

	log.Println(resp.Text)

	return resp.Text, nil
}

func (c Client) TextToFormattedRecipe(speechText string, descriptionText string) (string, error) {
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
		log.Printf("Chat completion error: %v\n", err)
		return "", err
	}

	log.Println(resp.Choices[0].Message.Content)

	return resp.Choices[0].Message.Content, nil
}
