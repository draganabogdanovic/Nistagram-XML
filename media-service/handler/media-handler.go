package handler

import (
	"net/http"

	"github.com/KristijanPill/Nishtagram/media-service/helpers"
	"github.com/KristijanPill/Nishtagram/media-service/payload"
	"github.com/KristijanPill/Nishtagram/media-service/service"
)

type MediaHandler struct {
	service *service.MediaService
}

func NewMediaHandler(service *service.MediaService) *MediaHandler {
	return &MediaHandler{service: service}
}

func (handler *MediaHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1024 * 128)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	formData := r.MultipartForm

	images := formData.File["media"]

	var postPaths []string

	for i, _ := range images {
		image, err := images[i].Open()
		defer image.Close()

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		path, err := handler.service.CreatePost(&image, images[i].Filename)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		postPaths = append(postPaths, path)
	}

	posts := &payload.MediaUploadResponse{MediaPaths: postPaths}
	helpers.ToJSON(&posts, w)
}

func (handler *MediaHandler) CreateStory(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1024 * 128)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	formData := r.MultipartForm

	images := formData.File["media"]

	var postPaths []string

	for i, _ := range images {
		image, err := images[i].Open()
		defer image.Close()

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		path, err := handler.service.CreateStory(&image, images[i].Filename)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		postPaths = append(postPaths, path)
	}

	posts := &payload.MediaUploadResponse{MediaPaths: postPaths}
	helpers.ToJSON(&posts, w)
}

func (handler *MediaHandler) CreateDocument(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1024 * 128)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	file, handle, err := r.FormFile("media")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	documentPath, err := handler.service.UploadDocumentPicture(&file, handle.Filename)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	documentPictureUploadResponse := &payload.DocumentPictureUploadResponse{DocumentPicture: documentPath}
	helpers.ToJSON(&documentPictureUploadResponse, w)
}

func (handler *MediaHandler) UploadProfilePicture(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1024 * 128)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	file, handle, err := r.FormFile("media")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	picturePath, err := handler.service.UploadProfilePicture(&file, handle.Filename)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	profilePictureUploadResponse := &payload.ProfilePictureUploadResponse{ProfilePicture: picturePath}
	helpers.ToJSON(&profilePictureUploadResponse, w)
}
