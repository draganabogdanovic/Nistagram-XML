package handler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/KristijanPill/Nishtagram/user-service/helpers"
	"github.com/KristijanPill/Nishtagram/user-service/middleware"
	"github.com/KristijanPill/Nishtagram/user-service/model"
	"github.com/KristijanPill/Nishtagram/user-service/payload"
	"github.com/KristijanPill/Nishtagram/user-service/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type VerificationRequestHandler struct {
	service *service.VerificationRequestService
}

func NewVerificationRequestHandler(service *service.VerificationRequestService) *VerificationRequestHandler {
	return &VerificationRequestHandler{service: service}
}

func (handler *VerificationRequestHandler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	r.Body = ioutil.NopCloser(bytes.NewReader(body))
	requestURL := fmt.Sprintf("http://%s:%s/upload/document", os.Getenv("MEDIA_SERVICE_DOMAIN"), os.Getenv("MEDIA_SERVICE_PORT"))
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

	documentPath := &payload.DocumentPictureUploadResponse{}
	err = helpers.FromJSON(&documentPath, response.Body)
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

		return
	}

	category, err := strconv.Atoi(r.FormValue("category"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	verReq := &model.VerificationRequest{
		UserID:                  userID,
		Name:                    r.FormValue("name"),
		Surname:                 r.FormValue("surname"),
		VerificationCategory:    model.VerificationCategory(category),
		OfficialDocumentPicture: documentPath.DocumentPicture,
	}

	_, err = handler.service.Create(verReq)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
}
