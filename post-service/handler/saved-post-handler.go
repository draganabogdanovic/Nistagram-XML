package handler

import (
	"net/http"

	"github.com/KristijanPill/Nishtagram/post-service/helpers"
	"github.com/KristijanPill/Nishtagram/post-service/middleware"
	"github.com/KristijanPill/Nishtagram/post-service/payload"
	"github.com/KristijanPill/Nishtagram/post-service/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type SavedPostHandler struct {
	service *service.SavedPostService
}

func NewSavedPostHandler(service *service.SavedPostService) *SavedPostHandler {
	return &SavedPostHandler{service: service}
}

func (handler *SavedPostHandler) SavePost(w http.ResponseWriter, r *http.Request) {
	dto := &payload.SavedPostCreate{}
	err := helpers.FromJSON(&dto, r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	claims := r.Context().Value(middleware.TokenKey{}).(jwt.MapClaims)

	userIDString := helpers.ExtractClaim("sub", claims)
	userID, err := uuid.Parse(userIDString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	dto.UserID = userID

	_, err = handler.service.SavePost(dto)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
}

func (handler *SavedPostHandler) GetAllCollectionNames(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.TokenKey{}).(jwt.MapClaims)

	userIDString := helpers.ExtractClaim("sub", claims)
	userID, err := uuid.Parse(userIDString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	collectionNames, err := handler.service.GetAllCollectionNames(userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	helpers.ToJSON(&collectionNames, w)
}

func (handler *SavedPostHandler) GetAllByLoggedInUser(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.TokenKey{}).(jwt.MapClaims)

	userIDString := helpers.ExtractClaim("sub", claims)
	userID, err := uuid.Parse(userIDString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	savedPosts, err := handler.service.GetAllByLoggedInUser(userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	helpers.ToJSON(&savedPosts, w)
}
