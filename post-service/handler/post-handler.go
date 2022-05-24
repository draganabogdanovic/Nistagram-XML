package handler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/KristijanPill/Nishtagram/post-service/helpers"
	"github.com/KristijanPill/Nishtagram/post-service/middleware"
	"github.com/KristijanPill/Nishtagram/post-service/model"
	"github.com/KristijanPill/Nishtagram/post-service/payload"
	"github.com/KristijanPill/Nishtagram/post-service/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type PostHandler struct {
	service *service.PostService
}

func NewPostHandler(service *service.PostService) *PostHandler {
	return &PostHandler{service: service}
}

func (handler *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	r.Body = ioutil.NopCloser(bytes.NewReader(body))
	requestURL := fmt.Sprintf("http://%s:%s/upload/post", os.Getenv("MEDIA_SERVICE_DOMAIN"), os.Getenv("MEDIA_SERVICE_PORT"))
	proxyReq, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(body))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	for header, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(header, value)
		}
	}

	client := &http.Client{}
	response, err := client.Do(proxyReq)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	postPaths := &payload.PostUploadResponse{}

	err = helpers.FromJSON(&postPaths, response.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	claims := r.Context().Value(middleware.TokenKey{}).(jwt.MapClaims)

	userIDString := helpers.ExtractClaim("sub", claims)
	userID, err := uuid.Parse(userIDString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	err = r.ParseMultipartForm(1024 * 128)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err.Error())
		return
	}

	formData := r.MultipartForm

	locationID, _ := uuid.Parse(r.FormValue("location_id"))

	post := &model.Post{
		UserID:      userID,
		Description: r.FormValue("description"),
		Tags:        formData.Value["tags"],
		Content:     postPaths.PostPaths,
		LocationID:  locationID,
	}

	_, err = handler.service.Create(post)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func (handler *PostHandler) FindByOtherUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDString := vars["id"]
	userID, err := uuid.Parse(userIDString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	var loggedInUserID uuid.UUID
	if r.Context().Value(middleware.LoggedInUser{}) != nil {
		loggedInUserID = r.Context().Value(middleware.LoggedInUser{}).(uuid.UUID)
	}

	posts, err := handler.service.FindByOtherUser(userID, loggedInUserID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	helpers.ToJSON(&posts, w)
}

func (handler *PostHandler) FindByUser(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.TokenKey{}).(jwt.MapClaims)

	userIDString := helpers.ExtractClaim("sub", claims)
	userID, err := uuid.Parse(userIDString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	posts, err := handler.service.FindByOtherUser(userID, userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	helpers.ToJSON(&posts, w)
}

func (handler *PostHandler) GetLikedPosts(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.TokenKey{}).(jwt.MapClaims)

	userIDString := helpers.ExtractClaim("sub", claims)
	userID, err := uuid.Parse(userIDString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	posts, err := handler.service.GetReviewsByUserIDAndStatus(userID, 1)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	helpers.ToJSON(&posts, w)

}

func (handler *PostHandler) GetDislikedPosts(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.TokenKey{}).(jwt.MapClaims)

	userIDString := helpers.ExtractClaim("sub", claims)
	userID, err := uuid.Parse(userIDString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	posts, err := handler.service.GetReviewsByUserIDAndStatus(userID, 0)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	helpers.ToJSON(&posts, w)

}

func (handler *PostHandler) CreateReport(w http.ResponseWriter, r *http.Request) {

	dto := &payload.ReportCreate{}
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

	_, err = handler.service.CreateReport(dto)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
}

func (handler *PostHandler) SearchPostsByLocation(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")

	var loggedInUserID uuid.UUID
	if r.Context().Value(middleware.LoggedInUser{}) != nil {
		loggedInUserID = r.Context().Value(middleware.LoggedInUser{}).(uuid.UUID)
	}

	posts, err := handler.service.SearchPostsByLocation(query, loggedInUserID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	helpers.ToJSON(&posts, w)
}

func (handler *PostHandler) SearchPostsByTags(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")

	var loggedInUserID uuid.UUID
	if r.Context().Value(middleware.LoggedInUser{}) != nil {
		loggedInUserID = r.Context().Value(middleware.LoggedInUser{}).(uuid.UUID)
	}

	posts, err := handler.service.SearchPostsByTags(query, loggedInUserID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	helpers.ToJSON(&posts, w)
}
