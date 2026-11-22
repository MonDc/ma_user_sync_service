package user

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mondc/ma_user_sync_service/internal/domain/errors"
)

type User struct {
	ID          string     `json:"id" db:"id"`
	PublicID    string     `json:"public_id" db:"public_id"`
	Email       string     `json:"email" db:"email"`
	FirstName   string     `json:"first_name" db:"first_name"`
	LastName    string     `json:"last_name" db:"last_name"`
	DateOfBirth time.Time  `json:"date_of_birth" db:"date_of_birth"`
	Status      UserStatus `json:"status" db:"status"`
	Metadata    JSON       `json:"metadata" db:"metadata"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	SyncedAt    *time.Time `json:"synced_at" db:"synced_at"`
}

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "suspended"
	UserStatusPending  UserStatus = "pending"
	UserStatusBlocked  UserStatus = "banned"
)

type JSON map[string]interface{}

// Scan implements the sql.Scanner interface to handle NULL values from the DB
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}
	return json.Unmarshal(b, j)
}

// Value implements the driver.Valuer interface for writing to the DB
func (j JSON) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Domain business rules
func (u *User) Validate() error {
	if u.Email == "" {
		return errors.ErrInvalidUserData(nil)
	}
	if u.FirstName == "" || u.LastName == "" {
		return errors.ErrInvalidUserData(nil)
	}
	if u.PublicID == "" {
		return errors.ErrInvalidUserData(nil)
	}
	if u.DateOfBirth.IsZero() {
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