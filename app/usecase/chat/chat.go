package chat

import (
	"context"
	"fmt"
	"log"
	"strings"

	redis_context "chatbot-service/app/domain/context"
	"chatbot-service/app/domain/persona"
	"chatbot-service/app/domain/session"
	"chatbot-service/dependencies/openai"
)

type ChatUseCase struct {
	PromptRepo    persona.PromptRepository
	ChatbotClient openai.Client
	SessionRepo   session.SessionRepository
	ContextRepo   redis_context.ContextRepository
}

func NewChatUseCase(
	promptRepo persona.PromptRepository,
	chatbot openai.Client,
	sessionRepo session.SessionRepository,
	contextRepo redis_context.ContextRepository,
) *ChatUseCase {
	return &ChatUseCase{
		PromptRepo:    promptRepo,
		ChatbotClient: chatbot,
		SessionRepo:   sessionRepo,
		ContextRepo:   contextRepo,
	}
}

func (uc *ChatUseCase) ProcessarMensagem(ctx context.Context, clienteID, phoneNumber, mensagem string) (string, error) {
    sessionID, err := uc.getOrCreateSessionID(ctx, phoneNumber)
    if err != nil {
        return "", err
    }

    chatCtx, err := uc.loadContext(ctx, sessionID)
    if err != nil {
        return "", err
    }

    systemPrompt, err := uc.buildSystemPrompt(ctx, clienteID, chatCtx)
    if err != nil {
        return "", err
    }

    mensagens := []openai.Message{
        {Role: "system", Content: systemPrompt},
        {Role: "user", Content: mensagem},
    }

    resposta, err := uc.ChatbotClient.Chat(ctx, mensagens)
    if err != nil {
        return "", fmt.Errorf("erro ao consultar o chatbot: %w", err)
    }

    err = uc.updateContext(ctx, sessionID, chatCtx, mensagem, resposta, phoneNumber)
    if err != nil {
        return "", err
    }

    return resposta, nil
}

func (uc *ChatUseCase) getOrCreateSessionID(ctx context.Context, phoneNumber string) (string, error) {
    sessionID, err := uc.SessionRepo.GetOrCreateSessionID(ctx, phoneNumber)
    if err != nil {
        return "", fmt.Errorf("erro ao buscar/criar sessionID: %w", err)
    }
    return sessionID, nil
}

func (uc *ChatUseCase) loadContext(ctx context.Context, sessionID string) (*redis_context.ChatContext, error) {
    chatCtx, err := uc.ContextRepo.LoadContext(ctx, sessionID)
    if err != nil {
        return nil, fmt.Errorf("erro ao buscar contexto: %w", err)
    }
    return &chatCtx, nil
}

func (uc *ChatUseCase) buildSystemPrompt(ctx context.Context, clienteID string, chatCtx *redis_context.ChatContext) (string, error) {
    prompt, err := uc.PromptRepo.GetPromptByClienteID(ctx, clienteID)
    if err != nil {
        return "", fmt.Errorf("erro ao buscar prompt: %w", err)
    }

    descricaoServicos := fmt.Sprintf("%s oferece os seguintes serviços:\n", prompt.NomeEmpresa)
    for _, s := range prompt.Servicos {
        descricaoServicos += fmt.Sprintf("- %s: %s — R$ %.2f\n", s.Nome, s.Descricao, s.Preco)
    }

    historico := ""
    for _, h := range chatCtx.History {
        historico += fmt.Sprintf("%s\n", h)
    }

    return fmt.Sprintf("%s\n\n%s\n\nHistórico recente:\n%s", prompt.SystemPrompt, descricaoServicos, historico), nil
}

func (uc *ChatUseCase) updateContext(ctx context.Context, sessionID string, chatCtx *redis_context.ChatContext, mensagem, resposta, phoneNumber string) error {
    chatCtx.History = append(chatCtx.History, "Usuário: "+mensagem)
    chatCtx.History = append(chatCtx.History, "Bot: "+resposta)

    if len(chatCtx.History) > 10 {
        chatCtx.History = chatCtx.History[len(chatCtx.History)-10:]
    }
	detectedIntent, err := uc.detectarIntentIA(ctx, mensagem, uc.intentsBasicas())
	if err != nil {
		log.Printf("Erro ao detectar intenção via IA, usando fallback: %v", err)
		detectedIntent = "Outra"
	}

	chatCtx.LastIntent = detectedIntent
	chatCtx.UserName = phoneNumber

    return uc.ContextRepo.SaveContext(ctx, sessionID, *chatCtx)
}

func (uc *ChatUseCase) intentsBasicas() []string {
    return []string{
        "Saudacao",
        "Despedida",
        "Perguntar preco",
        "Agendar servico",
        "Dúvida geral",
		"Reclamacaoo",
		"Feedback",
		"Promocaoo",
		"Cancelamento",
    }
}


func (uc *ChatUseCase) detectarIntentIA(ctx context.Context, mensagem string, intents []string) (string, error) {
    listaIntents := ""
    for _, intent := range intents {
        listaIntents += fmt.Sprintf("- %s\n", intent)
    }

    prompt := fmt.Sprintf(`Analise a seguinte mensagem e classifique em uma das intenções abaixo:
		%s
		Responda apenas com a intenção exata da lista. Se não encontrar uma correspondência clara, responda "Outra".
		Mensagem: "%s"`, listaIntents, mensagem)

    mensagens := []openai.Message{
        {Role: "system", Content: prompt},
    }

    resposta, err := uc.ChatbotClient.Chat(ctx, mensagens)
    if err != nil {
        return "", fmt.Errorf("erro ao detectar intenção via IA: %w", err)
    }

    resposta = strings.TrimSpace(resposta)

    for _, intent := range intents {
        if strings.EqualFold(resposta, intent) {
            return intent, nil
        }
    }

    return "Outra", nil
}