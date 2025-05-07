package bucket_test

import (
	"context"
	"testing"
	"chatbot-service/dependencies/bucket"
	"github.com/stretchr/testify/assert"
)

// MockBucket simula a leitura de um arquivo do bucket
type MockBucket struct {
	Content []byte
	Err     error
}

func (m *MockBucket) GetFile(ctx context.Context, path string) ([]byte, error) {
	return m.Content, m.Err
}

func TestGetPromptByClienteID(t *testing.T) {
	mockJSON := []byte(`{
		"nomeEmpresa": "Clínica Sorriso Ideal",
		"systemPrompt": "Você é Ana, a atendente virtual.",
		"servicos": [
			{ "nome": "Clareamento", "descricao": "Dentes brancos", "preco": 600 },
			{ "nome": "Limpeza", "descricao": "Remove tártaro", "preco": 200 }
		]
	}`)

	mock := &MockBucket{Content: mockJSON}
	repo := bucket.NewBucketPromptRepository(mock)

	prompt, err := repo.GetPromptByClienteID(context.Background(), "clinica-sorriso")
	assert.NoError(t, err)
	assert.Equal(t, "clinica-sorriso", prompt.ClienteID)
	assert.Equal(t, "Clínica Sorriso Ideal", prompt.NomeEmpresa)
	assert.Equal(t, "Você é Ana, a atendente virtual.", prompt.SystemPrompt)
	assert.Len(t, prompt.Servicos, 2)
	assert.Equal(t, "Clareamento", prompt.Servicos[0].Nome)
}