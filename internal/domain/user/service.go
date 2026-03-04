package user

import (
    "context"
    "time"
    "github.com/mondc/ma_user_sync_service/internal/domain/errors"
)

type Service interface {
    SyncUser(ctx context.Context, userID string) (*User, error)
    SyncAllUsers(ctx context.Context) ([]*User, error)
    GetUser(ctx context.Context, userID string) (*User, error)
}

type domainService struct {
    mainRepo Repository // mi_users
    localRepo Repository // ma_users
}

func NewDomainService(mainRepo, localRepo Repository) Service {
    return &domainService{
        mainRepo:  mainRepo,
        localRepo: localRepo,
    }
}

func (s *domainService) SyncUser(ctx context.Context, userID string) (*User, error) {
    // 1. Get user from main database (mi_users)
    user, err := s.mainRepo.FindByID(ctx, userID)
    if err != nil {
        return nil, errors.ErrUserNotFound(err)
    }

    // 2. Validate user data
    if err := user.Validate(); err != nil {
        return nil, err
    }

    // 3. Begin transaction for local DB
    tx, err := s.localRepo.BeginTx(ctx)
    if err != nil {
        return nil, errors.ErrSyncFailed(err)
    }

    // 4. Check if user exists locally
    exists, err := s.localRepo.Exists(ctx, user.ID)
    if err != nil {
        s.localRepo.RollbackTx(ctx, tx)
        return nil, errors.ErrSyncFailed(err)
    }

    // 5. Save or update local user
    user.UpdateSyncedAt()
    
    if exists {
        err = s.localRepo.Update(ctx, user)
    } else {
        err = s.localRepo.Save(ctx, user)
    }

    if err != nil {
        s.localRepo.RollbackTx(ctx, tx)
        return nil, errors.ErrSyncFailed(err)
    }

    // 6. Commit transaction
    if err := s.localRepo.CommitTx(ctx, tx); err != nil {
        return nil, errors.ErrSyncFailed(err)
    }

    return user, nil
}

func (s *domainService) SyncAllUsers(ctx context.Context) ([]*User, error) {
    var syncedUsers []*User
    offset := 0
    limit := 100

    for {
        // Batch fetch from main DB
        users, err := s.mainRepo.FindAll(ctx, limit, offset)
        if err != nil {
            return nil, errors.ErrSyncFailed(err)
        }

        if len(users) == 0 {
            break
        }

        // Sync each user
        for _, user := range users {
            syncedUser, err := s.SyncUser(ctx, user.ID)
            if err != nil {
                // Log error but continue with other users
                continue
            }
            syncedUsers = append(syncedUsers, syncedUser)
        }

        offset += limit
    }

    return syncedUsers, nil
}

func (s *domainService) GetUser(ctx context.Context, userID string) (*User, error) {
    user, err := s.localRepo.FindByID(ctx, userID)
    if err != nil {
        return nil, errors.ErrUserNotFound(err)
    }
    return user, nil
}