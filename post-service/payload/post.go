package payload

import (
	"github.com/KristijanPill/Nishtagram/post-service/model"
	"github.com/google/uuid"
)

type PostUploadResponse struct {
	PostPaths []string `json:"mediaPaths"`
}

type CommentCreate struct {
	UserID             uuid.UUID `json:"user_id"`
	PostID             uuid.UUID `json:"post_id"`
	Content            string    `json:"content"`
	RepliedToCommentID uuid.UUID `json:"replied_to_comment_id"`
}

type ReviewCreate struct {
	UserID       uuid.UUID          `json:"user_id"`
	PostID       uuid.UUID          `json:"post_id"`
	ReviewStatus model.ReviewStatus `json:"review_status"`
}

type ReportCreate struct {
	UserID      uuid.UUID `json:"user_id"`
	PostID      uuid.UUID `json:"post_id"`
	Description string    `json:"description"`
}

type SavedPostCreate struct {
	UserID         uuid.UUID `json:"user_id"`
	PostID         uuid.UUID `json:"post_id"`
	CollectionName string    `json:"collection_name"`
}

type PostView struct {
	ID               uuid.UUID      `json:"id"`
	UserID           uuid.UUID      `json:"user_id"`
	Username         string         `json:"username"`
	ProfilePicture   string         `json:"profile_picture,omitempty"`
	Content          []string       `json:"content"`
	NumberOfLikes    int64          `json:"number_of_likes"`
	NumberOfDislikes int64          `json:"number_of_dislikes"`
	NumberOfComments int64          `json:"number_of_comments"`
	Status           int            `json:"review_status"`
	Location         model.Location `json:"location,omitempty"`
	Description      string         `json:"description,omitempty"`
}

type SavedPostView struct {
	PostView       PostView `json:"post"`
	CollectionName string   `json:"collection_name"`
}

type CommentView struct {
	ID                 uuid.UUID `json:"id"`
	UserID             uuid.UUID `json:"user_id"`
	Username           string    `json:"username"`
	ProfilePicture     string    `json:"profile_picture"`
	Content            string    `json:"content"`
	RepliedToCommentID uuid.UUID `json:"replied_to_comment_id"`
}

type UserID struct {
	ID uuid.UUID `json:"id"`
}

type UserIDs struct {
	IDs []UserID `json:"ids"`
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
