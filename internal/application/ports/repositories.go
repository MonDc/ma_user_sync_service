package ports

import (
    "context"
    "github.com/mondc/ma_user_sync_service/internal/domain/user"
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