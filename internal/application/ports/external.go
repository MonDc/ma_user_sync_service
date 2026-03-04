package ports

import (
    "context"
    "github.com/mondc/ma_user_sync_service/internal/domain/user"
)

type UserSyncService interface {
    SyncUser(ctx context.Context, userID string) (*user.User, error)
    SyncAllUsers(ctx context.Context) ([]*user.User, error)
    GetUser(ctx context.Context, userID string) (*user.User, error)
}