package helpers

import (
	"crypto/rsa"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func VerifyToken(token string, publicKey *rsa.PublicKey) bool {
	parts := strings.Split(token, ".")
	err := jwt.SigningMethodRS256.Verify(strings.Join(parts[0:2], "."), parts[2], publicKey)

	return err == nil
}

func ExtractTokenFromHeader(header string) string {
	parts := strings.Split(header, " ")
	if len(parts) == 2 {
		return parts[1]
	}

	return ""
}

func ExtractClaim(key string, claims jwt.MapClaims) string {
	return claims[key].(string)
}
