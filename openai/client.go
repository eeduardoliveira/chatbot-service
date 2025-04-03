package openai

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"chatbot-service/db"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatPayload struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Choice struct {
	Message Message `json:"message"`
}

type ChatAPIResponse struct {
	Choices []Choice `json:"choices"`
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


func GetChatResponse(userMsg string) (string, error) {
	client := resty.New()

	payload := ChatPayload{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{Role: "user", Content: userMsg},
		},
	}

	var result ChatAPIResponse
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY")).
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		SetResult(&result).
		Post("https://api.openai.com/v1/chat/completions")

	if err != nil {
		return "", err
	}

	if resp.IsError() {
		var apiErr OpenAIError
		if jsonErr := json.Unmarshal(resp.Body(), &apiErr); jsonErr == nil {
			return "", fmt.Errorf("OpenAI error: %s (%s)", apiErr.Error.Message, apiErr.Error.Code)
		}
		return resp.String(), fmt.Errorf("Erro HTTP: %s", resp.String())
	}

	if len(result.Choices) == 0 {
		return resp.String(), fmt.Errorf("Resposta vazia da OpenAI: %s", resp.String())
	}

	return result.Choices[0].Message.Content, nil

}

func PerguntarAoChatbot(prompt string) (string, error) {
	client := resty.New()

	servicos, err := db.BuscarServicos()
if err != nil {
	return "", fmt.Errorf("Erro ao buscar serviços: %w", err)
}

descricaoServicos := "A clínica oferece os seguintes serviços:\n"
for _, s := range servicos {
	descricaoServicos += fmt.Sprintf("- %s: %s — R$ %.2f\n", s.Nome, s.Descricao, s.Preco)
}

systemPrompt := fmt.Sprintf(`Você é uma atendente virtual chamada Ana, da Clínica Sorriso Ideal.
Fale de forma simpática e acolhedora.

%s

A clínica atende das 08h às 18h, de segunda a sábado.

Você deve responder com base nessas informações. Se o paciente perguntar algo fora do escopo ou solicitar valores exatos e personalizados, oriente que será necessário entrar em contato com a recepção.`, descricaoServicos)

	req := ChatRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{Role: "system", Content: systemPrompt,
    },
			{Role: "user", Content: prompt},
		},
	}

	var respData ChatResponse

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY")).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetResult(&respData).
		Post("https://api.openai.com/v1/chat/completions")

	if err != nil {
		return "", err
	}

	if resp.IsError() || len(respData.Choices) == 0 {
		return "", fmt.Errorf("erro da OpenAI: %s", resp.String())
	}

	return respData.Choices[0].Message.Content, nil
}
