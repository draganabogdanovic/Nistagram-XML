package handler

import (
	"errors"
	"net/http"

	"github.com/KristijanPill/Nishtagram/auth-service/helpers"
	"github.com/KristijanPill/Nishtagram/auth-service/middleware"
	"github.com/KristijanPill/Nishtagram/auth-service/payload"
	"github.com/KristijanPill/Nishtagram/auth-service/service"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwk"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (handler *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	loginRequest := &payload.LoginRequest{}

	helpers.FromJSON(&loginRequest, r.Body)

	tokens, err := handler.service.Login(loginRequest)

	if err != nil {
		if errors.Is(err, service.ErrUnathorized) {
			http.Error(w, err.Error(), http.StatusUnauthorized)

			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	helpers.ToJSON(&tokens, w)
}

func (handler *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	credentials := &payload.Credentials{}
	helpers.FromJSON(&credentials, r.Body)

	err := handler.service.Register(credentials)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func (handler *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshTokenAndClaims := r.Context().Value(middleware.RefreshKey{}).(payload.RefreshTokenAndClaimsDTO)
	userIDString := helpers.ExtractClaim("sub", refreshTokenAndClaims.Claims)
	userID, err := uuid.Parse(userIDString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accessToken, err := handler.service.Refresh(refreshTokenAndClaims.RefreshToken, userID)

	if err != nil {
		if errors.Is(err, service.ErrUnathorized) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	helpers.ToJSON(&payload.RefreshResponse{AccessToken: accessToken}, w)
}

func (handler *AuthHandler) GetPublicKeys(w http.ResponseWriter, r *http.Request) {
	key := jwk.NewRSAPublicKey()
	key.FromRaw(handler.service.PublicKey)
	key.Set(jwk.KeyTypeKey, "RSA")

	var keys []jwk.Key

	keys = append(keys, key)

	jwks := &payload.JWKS{
		Keys: keys,
	}

	helpers.ToJSON(&jwks, w)
}
