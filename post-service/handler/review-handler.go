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

type ReviewHandler struct {
	service *service.ReviewService
}

func NewReviewHandler(service *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{service: service}
}

func (handler *ReviewHandler) ReviewPost(w http.ResponseWriter, r *http.Request) {
	dto := &payload.ReviewCreate{}
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

	_, err = handler.service.ReviewPost(dto)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
}
