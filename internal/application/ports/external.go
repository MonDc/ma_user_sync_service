package ports

import (
    "context"
    "github.com/yourusername/user-sync-service/internal/domain/user"
)

type UserSyncService interface {
    SyncUser(ctx context.Context, userID string) (*user.User, error)
    SyncAllUsers(ctx context.Context) ([]*user.User, error)
    GetUser(ctx context.Context, userID string) (*user.User, error)
}