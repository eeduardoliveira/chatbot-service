package openai

import (
	"context"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
)

// Message representa uma mensagem do chat
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest representa o payload enviado à API da OpenAI
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// ChatChoice representa uma resposta individual da OpenAI
type ChatChoice struct {
	Message Message `json:"message"`
}

// ChatResponse representa o formato de resposta da OpenAI
type ChatResponse struct {
	Choices []ChatChoice `json:"choices"`
}

// Client é a interface usada no domínio para abstrair o provedor de IA
type Client interface {
	Chat(ctx context.Context, messages []Message) (string, error)
}

// OpenAIClient é a implementação do client da OpenAI via HTTP
type OpenAIClient struct {
	apiKey string
	model  string
	url    string
	client *resty.Client
}

// NewClient cria um novo client da OpenAI com base nas variáveis de ambiente
func NewClient() *OpenAIClient {
	return &OpenAIClient{
		apiKey: os.Getenv("OPENAI_API_KEY"),
		model:  os.Getenv("OPENAI_MODEL"),
		url:    os.Getenv("OPENAI_API_URL"),
		client: resty.New(),
	}
}

// Chat envia a requisição para a API da OpenAI e retorna a resposta textual
func (c *OpenAIClient) Chat(ctx context.Context, messages []Message) (string, error) {
	req := ChatRequest{
		Model:    c.model,
		Messages: messages,
	}

	var respData ChatResponse

	resp, err := c.client.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+c.apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetResult(&respData).
		Post(c.url)

	if err != nil {
		return "", fmt.Errorf("erro na requisição OpenAI: %w", err)
	}

	if resp.IsError() || len(respData.Choices) == 0 {
		return "", fmt.Errorf("erro da OpenAI: %s", resp.String())
	}

	return respData.Choices[0].Message.Content, nil
}