package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Story struct {
	ID               uuid.UUID `gorm:"primary_key; unique; type:uuid;"`
	UserID           uuid.UUID
	CreatedAt        time.Time
	Content          pq.StringArray `gorm:"type:varchar(1000)[]"`
	CloseFriendsOnly bool
}

type Report struct {
	ID          uuid.UUID `gorm:"primary_key;  unique; type:uuid;"`
	UserID      uuid.UUID `gorm:"type:uuid"`
	StoryID     uuid.UUID `gorm:"foreign_key; not_unique; type:uuid;"`
	Description string
}

type StoryHighlight struct {
	UserID        uuid.UUID `gorm:"primaryKey; type:uuid"`
	StoryID       uuid.UUID `gorm:"primaryKey; type:uuid"`
	Story         Story
	CreatedAt     time.Time
	HighlightName string
}

func (s *Story) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.New()
	return nil
}

func (r *Report) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.New()
	return nil
}
