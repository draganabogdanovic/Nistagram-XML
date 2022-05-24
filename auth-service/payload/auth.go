package payload

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwk"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JWKS struct {
	Keys []jwk.Key `json:"keys"`
}

type Credentials struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
}

type RefreshToken struct {
	UserID uuid.UUID `gorm:"primary_key; type:uuid;"`
	Token  string
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}

type RefreshTokenAndClaimsDTO struct {
	RefreshToken string
	Claims       jwt.MapClaims
}
