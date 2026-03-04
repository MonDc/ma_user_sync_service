package user

import (
    "time"
    "github.com/mondc/ma_user_sync_service/internal/domain/errors"
)

type User struct {
    ID        string     `json:"id" db:"id"`
    Email     string     `json:"email" db:"email"`
    FirstName string     `json:"first_name" db:"first_name"`
    LastName  string     `json:"last_name" db:"last_name"`
    Status    UserStatus `json:"status" db:"status"`
    Metadata  JSON       `json:"metadata" db:"metadata"`
    CreatedAt time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
    SyncedAt  *time.Time `json:"synced_at" db:"synced_at"`
}

type UserStatus string

const (
    UserStatusActive   UserStatus = "ACTIVE"
    UserStatusInactive UserStatus = "INACTIVE"
    UserStatusPending  UserStatus = "PENDING"
    UserStatusBlocked  UserStatus = "BLOCKED"
)

type JSON map[string]interface{}

// Domain business rules
func (u *User) Validate() error {
    if u.Email == "" {
        return errors.ErrInvalidUserData(nil)
    }
    if u.FirstName == "" || u.LastName == "" {
        return errors.ErrInvalidUserData(nil)
    }
    return nil
}

func (u *User) Activate() {
    u.Status = UserStatusActive
    u.UpdatedAt = time.Now()
}

func (u *User) Deactivate() {
    u.Status = UserStatusInactive
    u.UpdatedAt = time.Now()
}

func (u *User) UpdateSyncedAt() {
    now := time.Now()
    u.SyncedAt = &now
    u.UpdatedAt = now
}