package mysql

import (
    "context"
    "database/sql"
    "fmt"
    
    "github.com/jmoiron/sqlx"
    "github.com/mondc/ma_user_sync_service/internal/domain/user"
    "github.com/mondc/ma_user_sync_service/internal/domain/errors"
    "go.uber.org/zap"
)

type localRepository struct {
    db     *sqlx.DB
    logger *zap.Logger
}

func NewLocalRepository(db *sqlx.DB, logger *zap.Logger) user.Repository {
    return &localRepository{
        db:     db,
        logger: logger,
    }
}

func (r *localRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
    query := `SELECT id, email, first_name, last_name, status, metadata, created_at, updated_at, synced_at 
              FROM ma_users WHERE id = ?`
    
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

func (r *localRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
    query := `SELECT id, email, first_name, last_name, status, metadata, created_at, updated_at, synced_at 
              FROM ma_users WHERE email = ?`
    
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

func (r *localRepository) FindAll(ctx context.Context, limit, offset int) ([]*user.User, error) {
    query := `SELECT id, email, first_name, last_name, status, metadata, created_at, updated_at, synced_at 
              FROM ma_users LIMIT ? OFFSET ?`
    
    var users []*user.User
    err := r.db.SelectContext(ctx, &users, query, limit, offset)
    if err != nil {
        r.logger.Error("failed to find all users", zap.Error(err))
        return nil, fmt.Errorf("failed to find users: %w", err)
    }
    
    return users, nil
}

func (r *localRepository) Save(ctx context.Context, user *user.User) error {
    query := `INSERT INTO ma_users (id, email, first_name, last_name, status, metadata, created_at, updated_at, synced_at)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
    
    _, err := r.db.ExecContext(ctx, query,
        user.ID, user.Email, user.FirstName, user.LastName,
        user.Status, user.Metadata, user.CreatedAt, user.UpdatedAt, user.SyncedAt)
    
    if err != nil {
        r.logger.Error("failed to save user", zap.Error(err), zap.String("user_id", user.ID))
        return fmt.Errorf("failed to save user: %w", err)
    }
    
    return nil
}

func (r *localRepository) Update(ctx context.Context, user *user.User) error {
    query := `UPDATE ma_users 
              SET email = ?, first_name = ?, last_name = ?, status = ?, 
                  metadata = ?, updated_at = ?, synced_at = ?
              WHERE id = ?`
    
    result, err := r.db.ExecContext(ctx, query,
        user.Email, user.FirstName, user.LastName,
        user.Status, user.Metadata, user.UpdatedAt, user.SyncedAt, user.ID)
    
    if err != nil {
        r.logger.Error("failed to update user", zap.Error(err), zap.String("user_id", user.ID))
        return fmt.Errorf("failed to update user: %w", err)
    }
    
    rows, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    
    if rows == 0 {
        return errors.ErrUserNotFound(nil)
    }
    
    return nil
}

func (r *localRepository) Delete(ctx context.Context, id string) error {
    query := `DELETE FROM ma_users WHERE id = ?`
    
    result, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        r.logger.Error("failed to delete user", zap.Error(err), zap.String("user_id", id))
        return fmt.Errorf("failed to delete user: %w", err)
    }
    
    rows, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    
    if rows == 0 {
        return errors.ErrUserNotFound(nil)
    }
    
    return nil
}

func (r *localRepository) Exists(ctx context.Context, id string) (bool, error) {
    query := `SELECT EXISTS(SELECT 1 FROM ma_users WHERE id = ?)`
    
    var exists bool
    err := r.db.GetContext(ctx, &exists, query, id)
    if err != nil {
        r.logger.Error("failed to check if user exists", zap.Error(err), zap.String("user_id", id))
        return false, fmt.Errorf("failed to check existence: %w", err)
    }
    
    return exists, nil
}

func (r *localRepository) BeginTx(ctx context.Context) (interface{}, error) {
    tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{
        Isolation: sql.LevelReadCommitted,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to begin transaction: %w", err)
    }
    return tx, nil
}

func (r *localRepository) CommitTx(ctx context.Context, tx interface{}) error {
    sqlTx, ok := tx.(*sqlx.Tx)
    if !ok {
        return fmt.Errorf("invalid transaction type")
    }
    
    if err := sqlTx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }
    return nil
}

func (r *localRepository) RollbackTx(ctx context.Context, tx interface{}) error {
    sqlTx, ok := tx.(*sqlx.Tx)
    if !ok {
        return fmt.Errorf("invalid transaction type")
    }
    
    if err := sqlTx.Rollback(); err != nil {
        return fmt.Errorf("failed to rollback transaction: %w", err)
    }
    return nil
}