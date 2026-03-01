package commands

import (
    "context"
    "github.com/yourusername/user-sync-service/internal/domain/user"
    "github.com/yourusername/user-sync-service/internal/application/ports"
)

type SyncUserCommand struct {
    UserID string
}

type SyncUserHandler struct {
    userService ports.UserSyncService
}

func NewSyncUserHandler(userService ports.UserSyncService) *SyncUserHandler {
    return &SyncUserHandler{
        userService: userService,
    }
}

func (h *SyncUserHandler) Handle(ctx context.Context, cmd SyncUserCommand) (*user.User, error) {
    return h.userService.SyncUser(ctx, cmd.UserID)
}