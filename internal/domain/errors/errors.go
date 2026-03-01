package errors

import "fmt"

type DomainError struct {
    Code    string
    Message string
    Err     error
}

func (e *DomainError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
    }
    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewDomainError(code, message string, err error) *DomainError {
    return &DomainError{
        Code:    code,
        Message: message,
        Err:     err,
    }
}

// Domain-specific errors
var (
    ErrUserNotFound      = func(err error) *DomainError { return NewDomainError("USER_NOT_FOUND", "user not found", err) }
    ErrUserAlreadyExists = func(err error) *DomainError { return NewDomainError("USER_ALREADY_EXISTS", "user already exists", err) }
    ErrInvalidUserData   = func(err error) *DomainError { return NewDomainError("INVALID_USER_DATA", "invalid user data", err) }
    ErrSyncFailed        = func(err error) *DomainError { return NewDomainError("SYNC_FAILED", "sync operation failed", err) }
)