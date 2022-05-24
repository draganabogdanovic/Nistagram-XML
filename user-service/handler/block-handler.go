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

type BlockHandler struct {
	service *service.BlockService
}

func NewBlockHandler(service *service.BlockService) *BlockHandler {
	return &BlockHandler{service: service}
}

func (handler *BlockHandler) Block(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	blockedUserIDString := vars["id"]

	blockedUserID, err := uuid.Parse(blockedUserIDString)

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

	err = handler.service.Block(userID, blockedUserID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
}

func (handler *BlockHandler) Unblock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	blockedUserIDString := vars["id"]

	blockedUserID, err := uuid.Parse(blockedUserIDString)

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

	err = handler.service.Unblock(userID, blockedUserID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
}
