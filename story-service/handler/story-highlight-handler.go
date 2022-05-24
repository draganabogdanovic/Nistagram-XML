package handler

import (
	"net/http"

	"github.com/KristijanPill/Nishtagram/story-service/helpers"
	"github.com/KristijanPill/Nishtagram/story-service/middleware"
	"github.com/KristijanPill/Nishtagram/story-service/payload"
	"github.com/KristijanPill/Nishtagram/story-service/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type StoryHighlightHandler struct {
	service *service.StoryHighlightService
}

func NewStoryHighlightHandler(service *service.StoryHighlightService) *StoryHighlightHandler {
	return &StoryHighlightHandler{service: service}
}

func (handler *StoryHighlightHandler) HighlightStory(w http.ResponseWriter, r *http.Request) {
	dto := &payload.StoryHighlightCreate{}
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

	_, err = handler.service.HighlightStory(dto)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
}

func (handler *StoryHighlightHandler) GetAllHighlightNames(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.TokenKey{}).(jwt.MapClaims)

	userIDString := helpers.ExtractClaim("sub", claims)
	userID, err := uuid.Parse(userIDString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	highlightNames, err := handler.service.GetAllHighlightNames(userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	helpers.ToJSON(&highlightNames, w)
}

func (handler *StoryHighlightHandler) GetAllByLoggedInUser(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.TokenKey{}).(jwt.MapClaims)

	userIDString := helpers.ExtractClaim("sub", claims)
	userID, err := uuid.Parse(userIDString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	highlights, err := handler.service.GetAllByLoggedInUser(userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	helpers.ToJSON(&highlights, w)
}
