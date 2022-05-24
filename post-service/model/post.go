package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Post struct {
	ID          uuid.UUID `gorm:"primary_key; unique; type:uuid;"`
	UserID      uuid.UUID `gorm:"type:uuid"`
	CreatedAt   time.Time
	Description string
	Content     pq.StringArray `gorm:"type:varchar(1000)[]"`
	Tags        pq.StringArray `gorm:"type:varchar(100)[]"`
	LocationID  uuid.UUID      `gorm:"foreign_key; not_unique; default:null"`
	Location    Location
}

type ReviewStatus int

const (
	DISLIKE ReviewStatus = iota
	LIKE
	NONE
)

type Review struct {
	PostID uuid.UUID `gorm:"primaryKey; type:uuid"`
	Post   Post
	UserID uuid.UUID `gorm:"primaryKey; type:uuid"`
	Status ReviewStatus
}

type Report struct {
	ID          uuid.UUID `gorm:"primaryKey; unique; type:uuid"`
	UserID      uuid.UUID `gorm:"type:uuid"`
	PostID      uuid.UUID `gorm:"foreign_key; not_unique; type:uuid;"`
	Description string
}

type Comment struct {
	ID                 uuid.UUID `gorm:"primaryKey; unique; type:uuid"`
	PostID             uuid.UUID `gorm:"type:uuid"`
	Post               Post
	UserID             uuid.UUID `gorm:"type:uuid"`
	Content            string
	RepliedToCommentID uuid.UUID `gorm:"type:uuid"`
}

type SavedPost struct {
	UserID         uuid.UUID `gorm:"primaryKey; type:uuid"`
	PostID         uuid.UUID `gorm:"primaryKey; type:uuid"`
	Post           Post
	CreatedAt      time.Time
	CollectionName string
}

type Location struct {
	ID      uuid.UUID `gorm:"primaryKey; unique; type:uuid" json:"id"`
	Country string    `json:"country"`
	City    string    `json:"city"`
}

func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return nil
}

func (c *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New()
	return nil
}

func (l *Location) BeforeCreate(tx *gorm.DB) (err error) {
	l.ID = uuid.New()
	return nil
}

func (r *Report) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.New()
	return nil
}
