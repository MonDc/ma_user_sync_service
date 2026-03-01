package queries

import (
    "context"
    "github.com/yourusername/user-sync-service/internal/domain/user"
    "github.com/yourusername/user-sync-service/internal/application/ports"
)

type GetUserQuery struct {
    UserID string
}

type GetUserHandler struct {
    userService ports.UserSyncService
}

func NewGetUserHandler(userService ports.UserSyncService) *GetUserHandler {
    return &GetUserHandler{
        userService: userService,
    }
}

func (h *GetUserHandler) Handle(ctx context.Context, query GetUserQuery) (*user.User, error) {
    return h.userService.GetUser(ctx, query.UserID)
}