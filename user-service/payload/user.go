package payload

import (
	"time"

	"github.com/KristijanPill/Nishtagram/user-service/model"
	"github.com/google/uuid"
)

type CreateUser struct {
	Username    string       `json:"username"`
	Password    string       `json:"password"`
	Email       string       `json:"email"`
	Name        string       `json:"name"`
	DOB         time.Time    `json:"dob"`
	Gender      model.Gender `json:"gender"`
	PhoneNumber string       `json:"phone_number"`
	Website     string       `json:"website,omitempty"`
	Bio         string       `json:"bio,omitempty"`
}

type Credentials struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
}

type UserInfo struct {
	Username       string       `json:"username,omitempty"`
	Email          string       `json:"email"`
	Name           string       `json:"name"`
	DOB            time.Time    `json:"dob"`
	Gender         model.Gender `json:"gender"`
	PhoneNumber    string       `json:"phone_number"`
	Website        string       `json:"website"`
	Bio            string       `json:"bio"`
	ProfilePicture string       `json:"profile_picture,omitempty"`
}

type UserProfileInfo struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	Name           string    `json:"name"`
	Bio            string    `json:"bio"`
	Website        string    `json:"website"`
	Followers      int64     `json:"followers_count"`
	Following      int64     `json:"following_count"`
	ProfilePicture string    `json:"profile_picture,omitempty"`
}

type Username struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

type UserID struct {
	ID uuid.UUID `json:"id"`
}

type UserIDs struct {
	IDs []UserID `json:"ids"`
}

type Usernames struct {
	Usernames []Username `json:"usernames"`
}

type UsersDetails struct {
	UsersDetails []UserDetails `json:"users_details"`
}

type UserDetails struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	Private        bool      `json:"private"`
	Followed       bool      `json:"followed"`
	ProfilePicture string    `json:"profile_picture"`
}

type CreateVerificationRequest struct {
	ID                      uuid.UUID                  `json:"id"`
	UserID                  uuid.UUID                  `json:"user_id"`
	Name                    string                     `json:"name"`
	Surname                 string                     `json:"surname"`
	Verified                bool                       `json:"verified"`
	OfficialDocumentPicture string                     `json:"officialDocumentPicture"`
	VerificationCategory    model.VerificationCategory `json:"verificationCategory"`
}

type DocumentPictureUploadResponse struct {
	DocumentPicture string `json:"document_picture"`
}

type FollowStatus struct {
	IsFollower  bool `json:"is_follower"`
	CloseFriend bool `json:"close_friend"`
}

type ProfilePictureUploadResponse struct {
	ProfilePicture string `json:"profile_picture"`
}

type UserView struct {
	Username       string `json:"username"`
	Name           string `json:"name"`
	ProfilePicture string `json:"profile_picture,omitempty"`
}
