package handler

import (
	"net/http"

	"github.com/KristijanPill/Nishtagram/user-service/helpers"
	"github.com/KristijanPill/Nishtagram/user-service/middleware"
	"github.com/KristijanPill/Nishtagram/user-service/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type FollowRequestHandler struct {
	service *service.FollowRequestService
}

func NewFollowRequestHandler(service *service.FollowRequestService) *FollowRequestHandler {
	return &FollowRequestHandler{service: service}
}

func (handler *FollowRequestHandler) Accept(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	followerIDString := vars["id"]
	followerID, err := uuid.Parse(followerIDString)

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

	err = handler.service.Accept(userID, followerID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
}

func (handler *FollowRequestHandler) Decline(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	followerIDString := vars["id"]
	followerID, err := uuid.Parse(followerIDString)

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

	err = handler.service.Decline(userID, followerID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
}
