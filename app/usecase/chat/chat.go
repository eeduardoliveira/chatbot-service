package chat

import (
	"context"
	"fmt"

	"chatbot-service/app/domain/persona"
	"chatbot-service/dependencies/openai"
)

type ChatUseCase struct {
	PromptRepo    persona.PromptRepository
	ChatbotClient openai.Client
}

func NewChatUseCase(promptRepo persona.PromptRepository, chatbot openai.Client) *ChatUseCase {
	return &ChatUseCase{
		PromptRepo:    promptRepo,
		ChatbotClient: chatbot,
	}
}

func (uc *ChatUseCase) ProcessarMensagem(ctx context.Context, clienteID, mensagem string) (string, error) {
	prompt, err := uc.PromptRepo.GetPromptByClienteID(ctx, clienteID)
	if err != nil {
		return "", fmt.Errorf("erro ao buscar prompt: %w", err)
	}

	descricaoServicos := fmt.Sprintf("%s oferece os seguintes serviços:\n", prompt.NomeEmpresa)
	for _, s := range prompt.Servicos {
		descricaoServicos += fmt.Sprintf("- %s: %s — R$ %.2f\n", s.Nome, s.Descricao, s.Preco)
	}

	systemPrompt := fmt.Sprintf("%s\n\n%s", prompt.SystemPrompt, descricaoServicos)

	mensagens := []openai.Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: mensagem},
	}

	resposta, err := uc.ChatbotClient.Chat(ctx, mensagens)
	if err != nil {
		return "", fmt.Errorf("erro ao consultar o chatbot: %w", err)
	}

	return resposta, nil
}