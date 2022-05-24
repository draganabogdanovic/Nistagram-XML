package payload

import (
	"github.com/google/uuid"
)

type LoginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	IsLoginSuccessful bool      `json:"is_login_successful"`
	ID                uuid.UUID `json:"id"`
	Error             string    `json:"error"`
}
