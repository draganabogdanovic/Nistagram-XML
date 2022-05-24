package middleware

import (
	"context"
	"crypto/rsa"
	"net/http"

	"github.com/KristijanPill/Nishtagram/auth-service/helpers"
	"github.com/KristijanPill/Nishtagram/auth-service/payload"
	"github.com/dgrijalva/jwt-go"
)

type RefreshMiddleware struct {
	publicKey *rsa.PublicKey
	hmacKey   []byte
}

func NewRefreshMiddleware(publicKey *rsa.PublicKey, hmacKey []byte) *RefreshMiddleware {
	return &RefreshMiddleware{
		publicKey: publicKey,
		hmacKey:   hmacKey,
	}
}

type RefreshKey struct{}

func (middleware *RefreshMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Authorization"] == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		accessTokenString := helpers.ExtractTokenFromHeader(r.Header["Authorization"][0])

		accessTokenClaims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(accessTokenString, &accessTokenClaims, func(token *jwt.Token) (interface{}, error) {
			return middleware.publicKey, nil
		})

		v, _ := err.(*jwt.ValidationError)
		if v.Errors != jwt.ValidationErrorExpired {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		refreshRequest := &payload.RefreshRequest{}
		helpers.FromJSON(&refreshRequest, r.Body)

		refreshTokenClaims := jwt.MapClaims{}
		_, err = jwt.ParseWithClaims(refreshRequest.RefreshToken, &refreshTokenClaims, func(token *jwt.Token) (interface{}, error) {
			return middleware.hmacKey, nil
		})

		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), RefreshKey{}, payload.RefreshTokenAndClaimsDTO{
			RefreshToken: refreshRequest.RefreshToken,
			Claims:       refreshTokenClaims,
		})

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
