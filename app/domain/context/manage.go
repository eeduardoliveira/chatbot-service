package context

import "context"

// UseCase encapsula o caso de uso para gerenciar contexto de conversa.
type UseCase struct {
    Repository ContextRepository
}

// NewUseCase cria uma nova instância do caso de uso de contexto.
func NewUseCase(repo ContextRepository) *UseCase {
    return &UseCase{
        Repository: repo,
    }
}

// LoadContext busca o contexto no Redis (ou retorna contexto novo se não existir).
func (u *UseCase) LoadContext(ctx context.Context, sessionID string) (ChatContext, error) {
    return u.Repository.LoadContext(ctx, sessionID)
}

// SaveContext atualiza o contexto no Redis.
func (u *UseCase) SaveContext(ctx context.Context, sessionID string, chatCtx ChatContext) error {
    return u.Repository.SaveContext(ctx, sessionID, chatCtx)
}