package handler

import (
	"net/http"

	"github.com/KristijanPill/Nishtagram/post-service/helpers"
	"github.com/KristijanPill/Nishtagram/post-service/middleware"
	"github.com/KristijanPill/Nishtagram/post-service/payload"
	"github.com/KristijanPill/Nishtagram/post-service/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type CommentHandler struct {
	service *service.CommentService
}

func NewCommentHandler(service *service.CommentService) *CommentHandler {
	return &CommentHandler{service: service}
}

func (handler *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	dto := &payload.CommentCreate{}
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

	_, err = handler.service.Create(dto)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
}

func (handler *CommentHandler) FindAllByPostID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	postIDString := vars["id"]
	postID, err := uuid.Parse(postIDString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	comments, err := handler.service.FindAllByPostID(postID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	helpers.ToJSON(comments, w)
}
