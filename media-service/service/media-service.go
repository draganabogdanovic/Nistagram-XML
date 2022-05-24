package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"

	"github.com/google/uuid"
)

type MediaService struct{}

const (
	postPath           = "/storage/posts/"
	storyPath          = "/storage/stories/"
	profilePicturePath = "/storage/profile_pictures/"
	documentPicturePath = "/storage/documents/"
)

func NewMediaService() *MediaService {
	return &MediaService{}
}

func (service *MediaService) CreatePost(image *multipart.File, filename string) (string, error) {
	mediaName := uuid.New().String() + "." + service.extractFileExtension(filename)

	destination, err := os.Create("." + postPath + mediaName)
	defer destination.Close()

	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	if _, err := io.Copy(destination, *image); err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return postPath + mediaName, nil
}

func (service *MediaService) CreateStory(image *multipart.File, filename string) (string, error) {
	mediaName := uuid.New().String() + "." + service.extractFileExtension(filename)

	destination, err := os.Create("." + storyPath + mediaName)
	defer destination.Close()

	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	if _, err := io.Copy(destination, *image); err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return storyPath + mediaName, nil
}

func (service *MediaService) UploadProfilePicture(image *multipart.File, filename string) (string, error) {
	mediaName := uuid.New().String() + "." + service.extractFileExtension(filename)

	destination, err := os.Create("." + profilePicturePath + mediaName)
	defer destination.Close()

	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	if _, err := io.Copy(destination, *image); err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return profilePicturePath + mediaName, nil
}

func (*MediaService) extractFileExtension(filename string) string {
	parts := strings.Split(filename, ".")
	return parts[len(parts)-1]
}

func (service *MediaService) UploadDocumentPicture(image *multipart.File, filename string) (string, error) {
	mediaName := uuid.New().String() + "." + service.extractFileExtension(filename)

	destination, err := os.Create("." + documentPicturePath + mediaName)
	defer destination.Close()

	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	if _, err := io.Copy(destination, *image); err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return documentPicturePath + mediaName, nil
}
