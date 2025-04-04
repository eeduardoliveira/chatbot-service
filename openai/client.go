package openai

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIError struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ChatChoice struct {
	Message Message `json:"message"`
}

type ChatResponse struct {
	Choices []ChatChoice `json:"choices"`
}


func PerguntarAoChatbot(prompt string, historico []Message) ([]Message, string, error) {
	client := resty.New()


historico = append(historico, Message{Role: "user", Content: prompt})


req := ChatRequest{
	Model:    "gpt-4o",
	Messages: historico,
}

	var respData ChatResponse

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY")).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetResult(&respData).
		Post("https://api.openai.com/v1/chat/completions")

	if err != nil {
		return nil, "", err
	}

	if resp.IsError() || len(respData.Choices) == 0 {
		return nil, "", fmt.Errorf("erro da OpenAI: %s", resp.String())
	}

	return nil, respData.Choices[0].Message.Content, err
}
