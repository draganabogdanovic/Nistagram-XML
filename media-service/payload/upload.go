package payload

type MediaUploadResponse struct {
	MediaPaths []string `json:"mediaPaths"`
}

type DocumentPictureUploadResponse struct {
	DocumentPicture string `json:"document_picture"`
}

type ProfilePictureUploadResponse struct {
	ProfilePicture string `json:"profile_picture"`
}
