package bucket

import (
	"context"
	"encoding/json"
	"fmt"
	"chatbot-service/app/domain/persona"
)


type BucketPromptRepository struct {
	Bucket BucketService
}

// NewBucketPromptRepository cria um novo repositório de prompt baseado em bucket
func NewBucketPromptRepository(bucket BucketService) *BucketPromptRepository {
	return &BucketPromptRepository{
		Bucket: bucket,
	}
}

func (r *BucketPromptRepository) GetPromptByClienteID(ctx context.Context, clienteID string) (*persona.Prompt, error) {
	body, err := r.Bucket.GetFile(ctx, clienteID)
	if err != nil {
		return nil, err
	}

	var data struct {
		NomeEmpresa  string              `json:"nomeEmpresa"`
		SystemPrompt string              `json:"systemPrompt"`
		Servicos     []persona.Servico   `json:"servicos"`
	}
	
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("erro ao fazer unmarshal do JSON: %w", err)
	}
	
	return &persona.Prompt{
		ClienteID:    clienteID,
		NomeEmpresa:  data.NomeEmpresa,
		SystemPrompt: data.SystemPrompt,
		Servicos:     data.Servicos,
	}, nil
}