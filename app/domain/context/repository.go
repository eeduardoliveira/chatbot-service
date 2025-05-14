package context
import "context"

type ChatContext struct {
    History     []string `json:"history"`
    LastIntent  string   `json:"last_intent"`
    UserName    string   `json:"user_name"`
    BusinessType string  `json:"business_type"`
}

type ContextRepository interface {
    LoadContext(ctx context.Context, sessionID string) (ChatContext, error)
    SaveContext(ctx context.Context, sessionID string, context ChatContext) error
}