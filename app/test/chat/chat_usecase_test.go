package chat_test

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"chatbot-service/app/domain/persona"
	"chatbot-service/app/usecase/chat"
	"chatbot-service/dependencies/openai"

	"github.com/stretchr/testify/assert"
)

// Mock do PromptRepository
type mockPromptRepo struct {
	Prompt *persona.Prompt
	Err    error
}

func (m *mockPromptRepo) GetPromptByClienteID(ctx context.Context, clienteID string) (*persona.Prompt, error) {
	return m.Prompt, m.Err
}

// Mock do ChatbotClient
type mockChatbotClient struct {
	Response string
	Err      error
}

func (m *mockChatbotClient) Chat(ctx context.Context, messages []openai.Message) (string, error) {
	return m.Response, m.Err
}

func loadMockPrompt(t *testing.T) *persona.Prompt {
	t.Helper()

	filePath := filepath.Join("..", "testdata", "prompt_client.json")
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Erro ao ler mock prompt: %v", err)
	}

	var prompt persona.Prompt
	if err := json.Unmarshal(content, &prompt); err != nil {
		t.Fatalf("Erro ao fazer unmarshal do prompt: %v", err)
	}

	return &prompt
}
func TestChatUseCase_Sucesso(t *testing.T) {
	mockPrompt := loadMockPrompt(t)

	repo := &mockPromptRepo{Prompt: mockPrompt}
	client := &mockChatbotClient{Response: "Olá, como posso ajudar?"}

	useCase := chat.NewChatUseCase(repo, client)

	resp, err := useCase.ProcessarMensagem(context.Background(), "clinica-sorriso", "Oi")
	assert.NoError(t, err)
	assert.Equal(t, "Olá, como posso ajudar?", resp)
}

func TestChatUseCase_ErroPrompt(t *testing.T) {
	repo := &mockPromptRepo{Err: errors.New("erro ao buscar prompt")}
	client := &mockChatbotClient{}
	useCase := chat.NewChatUseCase(repo, client)

	resp, err := useCase.ProcessarMensagem(context.Background(), "cliente-z", "Oi")
	assert.Error(t, err)
	assert.Empty(t, resp)
}

func TestChatUseCase_ErroChatbot(t *testing.T) {
	mockPrompt := &persona.Prompt{
		ClienteID:    "clinica-x",
		NomeEmpresa:  "Clínica X",
		SystemPrompt: "Você é um assistente virtual.",
		Servicos:     []persona.Servico{},
	}

	repo := &mockPromptRepo{Prompt: mockPrompt}
	client := &mockChatbotClient{Err: errors.New("erro na IA")}
	useCase := chat.NewChatUseCase(repo, client)

	resp, err := useCase.ProcessarMensagem(context.Background(), "clinica-x", "Oi")
	assert.Error(t, err)
	assert.Empty(t, resp)
}