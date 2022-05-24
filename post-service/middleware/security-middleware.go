package middleware

import (
	"context"
	"crypto/rsa"
	"net/http"

	"github.com/KristijanPill/Nishtagram/post-service/helpers"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type SecurityMiddleware struct {
	PublicKey *rsa.PublicKey
}

func NewSecurityMiddleware(publicKey *rsa.PublicKey) *SecurityMiddleware {
	return &SecurityMiddleware{PublicKey: publicKey}
}

type TokenKey struct{}

func (middleware *SecurityMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Authorization"] == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)

			return
		}

		tokenString := helpers.ExtractTokenFromHeader(r.Header["Authorization"][0])

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return middleware.PublicKey, nil
		})

		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)

			return
		}

		if !token.Valid {
			http.Error(w, "unauthorized", http.StatusUnauthorized)

			return
		}

		ctx := context.WithValue(r.Context(), TokenKey{}, claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

type LoggedInUser struct{}

func (middleware *SecurityMiddleware) UserContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Authorization"] == nil {
			ctx := context.WithValue(r.Context(), LoggedInUser{}, nil)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
			return
		}

		tokenString := helpers.ExtractTokenFromHeader(r.Header["Authorization"][0])

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return middleware.PublicKey, nil
		})

		if err != nil {
			ctx := context.WithValue(r.Context(), LoggedInUser{}, nil)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
			return
		}

		if !token.Valid {
			ctx := context.WithValue(r.Context(), LoggedInUser{}, nil)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
			return
		}

		userIDString := helpers.ExtractClaim("sub", claims)
		userID, err := uuid.Parse(userIDString)

		if err != nil {
			ctx := context.WithValue(r.Context(), LoggedInUser{}, nil)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), LoggedInUser{}, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
