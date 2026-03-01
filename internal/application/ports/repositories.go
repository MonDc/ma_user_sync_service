package ports

import (
    "context"
    "github.com/yourusername/user-sync-service/internal/domain/user"
)

type UserRepository interface {
    user.Repository
}

type MainUserRepository interface {
    user.Repository
}

type LocalUserRepository interface {
    user.Repository
}