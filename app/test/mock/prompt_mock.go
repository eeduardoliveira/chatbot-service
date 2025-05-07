package mocks

import (
	"context"

	"chatbot-service/app/domain/persona"
)

type PromptMock struct {
	Response *persona.Prompt
	Err      error
}

func (m *PromptMock) GetPromptByClienteID(ctx context.Context, clienteID string) (*persona.Prompt, error) {
	return m.Response, m.Err
}