package handler

import (
	"errors"
	"net/http"

	"github.com/KristijanPill/Nishtagram/user-service/helpers"
	"github.com/KristijanPill/Nishtagram/user-service/middleware"
	"github.com/KristijanPill/Nishtagram/user-service/payload"
	"github.com/KristijanPill/Nishtagram/user-service/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type FollowHandler struct {
	service *service.FollowService
}

func NewFollowHandler(service *service.FollowService) *FollowHandler {
	return &FollowHandler{service: service}
}

func (handler *FollowHandler) Follow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	followedUserIDString := vars["id"]

	followedUserID, err := uuid.Parse(followedUserIDString)

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

	_, err = handler.service.Follow(userID, followedUserID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
}

func (handler *FollowHandler) AddCloseFriend(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	closeFriendIDString := vars["id"]
	closeFriendID, err := uuid.Parse(closeFriendIDString)

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

	err = handler.service.UpdateCloseFriend(closeFriendID, userID, true)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
}

func (handler *FollowHandler) RemoveCloseFriend(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	closeFriendIDString := vars["id"]
	closeFriendID, err := uuid.Parse(closeFriendIDString)

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

	err = handler.service.UpdateCloseFriend(closeFriendID, userID, false)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
}

func (handler *FollowHandler) GetOutgoingFollowStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	followedIDString := vars["id"]
	followedID, err := uuid.Parse(followedIDString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	var loggedInUserID uuid.UUID
	if r.Context().Value(middleware.LoggedInUser{}) != nil {
		loggedInUserID = r.Context().Value(middleware.LoggedInUser{}).(uuid.UUID)
	}

	follow, err := handler.service.FindByUserIDAndFollowerID(followedID, loggedInUserID)

	var response *payload.FollowStatus

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response = &payload.FollowStatus{
				IsFollower:  false,
				CloseFriend: false,
			}
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}
	} else {
		response = &payload.FollowStatus{
			IsFollower:  true,
			CloseFriend: follow.CloseFriend,
		}
	}

	helpers.ToJSON(&response, w)
}

func (handler *FollowHandler) Mute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	mutedIDString := vars["id"]
	mutedID, err := uuid.Parse(mutedIDString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	claims := r.Context().Value(middleware.TokenKey{}).(jwt.MapClaims)

	loggedInUserIDString := helpers.ExtractClaim("sub", claims)
	loggedInUserID, err := uuid.Parse(loggedInUserIDString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	err = handler.service.UpdateMuted(mutedID, loggedInUserID, true)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
}

func (handler *FollowHandler) Unmute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	mutedIDString := vars["id"]
	mutedID, err := uuid.Parse(mutedIDString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	claims := r.Context().Value(middleware.TokenKey{}).(jwt.MapClaims)

	loggedInUserIDString := helpers.ExtractClaim("sub", claims)
	loggedInUserID, err := uuid.Parse(loggedInUserIDString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	err = handler.service.UpdateMuted(mutedID, loggedInUserID, false)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
}
