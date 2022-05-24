package payload

import (
	"time"

	"github.com/google/uuid"
)

type StoryUploadResponse struct {
	StoryPaths []string `json:"mediaPaths"`
}

type StoryHighlightCreate struct {
	UserID        uuid.UUID `json:"user_id"`
	StoryID       uuid.UUID `json:"story_id"`
	HighlightName string    `json:"highlight_name"`
}

type ReportCreate struct {
	UserID      uuid.UUID `json:"user_id"`
	StoryID     uuid.UUID `json:"story_id"`
	Description string    `json:"description"`
}

type StoryView struct {
	ID               uuid.UUID `json:"id"`
	CreatedAt        time.Time `json:"created_at"`
	Content          []string  `json:"content"`
	CloseFriendsOnly bool      `json:"close_friends_only"`
}

type FollowStatus struct {
	IsFollower  bool `json:"is_follower"`
	CloseFriend bool `json:"close_friend"`
}

type StoryHighlightView struct {
	StoryView     StoryView `json:"story"`
	HighlightName string    `json:"highlight_name"`
}
