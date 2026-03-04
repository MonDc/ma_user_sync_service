package mysql

import (
    "context"
    "database/sql"
    "fmt"
    
    "github.com/jmoiron/sqlx"
    _ "github.com/go-sql-driver/mysql"
    "github.com/mondc/ma_user_sync_service/internal/domain/user"
    "github.com/mondc/ma_user_sync_service/internal/domain/errors"
    "go.uber.org/zap"
)

type mainRepository struct {
    db     *sqlx.DB
    logger *zap.Logger
}

func NewMainRepository(db *sqlx.DB, logger *zap.Logger) user.Repository {
    return &mainRepository{
        db:     db,
        logger: logger,
    }
}

func (r *mainRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
    query := `SELECT id, email, first_name, last_name, status, metadata, created_at, updated_at 
              FROM mi_users WHERE id = ?`
    
    var u user.User
    err := r.db.GetContext(ctx, &u, query, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.ErrUserNotFound(err)
        }
        r.logger.Error("failed to find user by ID", zap.Error(err), zap.String("user_id", id))
        return nil, fmt.Errorf("failed to find user: %w", err)
    }
    
    return &u, nil
}

func (r *mainRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
    query := `SELECT id, email, first_name, last_name, status, metadata, created_at, updated_at 
              FROM mi_users WHERE email = ?`
    
    var u user.User
    err := r.db.GetContext(ctx, &u, query, email)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.ErrUserNotFound(err)
        }
        r.logger.Error("failed to find user by email", zap.Error(err), zap.String("email", email))
        return nil, fmt.Errorf("failed to find user: %w", err)
    }
    
    return &u, nil
}

func (r *mainRepository) FindAll(ctx context.Context, limit, offset int) ([]*user.User, error) {
    query := `SELECT id, email, first_name, last_name, status, metadata, created_at, updated_at 
              FROM mi_users LIMIT ? OFFSET ?`
    
    var users []*user.User
    err := r.db.SelectContext(ctx, &users, query, limit, offset)
    if err != nil {
        r.logger.Error("failed to find all users", zap.Error(err))
        return nil, fmt.Errorf("failed to find users: %w", err)
    }
    
    return users, nil
}

func (r *mainRepository) Save(ctx context.Context, user *user.User) error {
    // Main repository is read-only for this service
    return errors.NewDomainError("READ_ONLY", "main repository is read-only", nil)
}

func (r *mainRepository) Update(ctx context.Context, user *user.User) error {
    // Main repository is read-only for this service
    return errors.NewDomainError("READ_ONLY", "main repository is read-only", nil)
}

func (r *mainRepository) Delete(ctx context.Context, id string) error {
    // Main repository is read-only for this service
    return errors.NewDomainError("READ_ONLY", "main repository is read-only", nil)
}

func (r *mainRepository) Exists(ctx context.Context, id string) (bool, error) {
    query := `SELECT EXISTS(SELECT 1 FROM mi_users WHERE id = ?)`
    
    var exists bool
    err := r.db.GetContext(ctx, &exists, query, id)
    if err != nil {
        r.logger.Error("failed to check if user exists", zap.Error(err), zap.String("user_id", id))
        return false, fmt.Errorf("failed to check existence: %w", err)
    }
    
    return exists, nil
}

func (r *mainRepository) BeginTx(ctx context.Context) (interface{}, error) {
    // Main repository is read-only, no transactions needed
    return nil, nil
}

func (r *mainRepository) CommitTx(ctx context.Context, tx interface{}) error {
    return nil
}

func (r *mainRepository) RollbackTx(ctx context.Context, tx interface{}) error {
    return nil
}