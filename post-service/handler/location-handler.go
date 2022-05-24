package handler

import (
	"net/http"

	"github.com/KristijanPill/Nishtagram/post-service/helpers"
	"github.com/KristijanPill/Nishtagram/post-service/service"
	"github.com/gorilla/mux"
)

type LocationHandler struct {
	service *service.LocationService
}

func NewLocationHandler(service *service.LocationService) *LocationHandler {
	return &LocationHandler{service: service}
}

func (handler *LocationHandler) GetByQuery(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	query := vars["query"]

	locations, err := handler.service.GetByQuery(query)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	helpers.ToJSON(&locations, w)
}
