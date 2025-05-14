package session

import "context"

type SessionRepository interface {
    GetSessionID(ctx context.Context, phoneNumber string) (string, error)
    SaveSessionID(ctx context.Context, phoneNumber, sessionID string) error
    GetOrCreateSessionID(ctx context.Context, phoneNumber string) (string, error)
}