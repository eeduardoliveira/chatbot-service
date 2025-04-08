package persona

import "context"

// Prompt representa o prompt dinâmico carregado para cada cliente
type Servico struct {
	Nome      string  `json:"nome"`
	Descricao string  `json:"descricao"`
	Preco     float64 `json:"preco"`
}

type Prompt struct {
	ClienteID    string
	NomeEmpresa  string     `json:"nomeEmpresa"`
	SystemPrompt string     `json:"systemPrompt"`
	Servicos     []Servico  `json:"servicos"`
}

// PromptRepository define o contrato para carregar o prompt customizado para um cliente
type PromptRepository interface {
	GetPromptByClienteID(ctx context.Context, clienteID string) (*Prompt, error)
}