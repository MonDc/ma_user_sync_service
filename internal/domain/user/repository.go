package user

import (
    "context"
    "github.com/mondc/ma_user_sync_service/internal/domain/errors"
)

type Repository interface {
    // Main repository operations (mi_users)
    FindByID(ctx context.Context, id string) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)
    FindAll(ctx context.Context, limit, offset int) ([]*User, error)
    
    // Local repository operations (ma_users)
    Save(ctx context.Context, user *User) error
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id string) error
    Exists(ctx context.Context, id string) (bool, error)
    
    // Sync operations
    BeginTx(ctx context.Context) (interface{}, error)
    CommitTx(ctx context.Context, tx interface{}) error
    RollbackTx(ctx context.Context, tx interface{}) error
}