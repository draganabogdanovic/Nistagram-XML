package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID                     uuid.UUID `gorm:"primary_key; unique; type:uuid;"`
	Role                   Role
	Username               string
	Email                  string
	ProfilePicture         string
	Private                bool
	Verified               bool
	Taggable               bool
	CanRecieveAnonMessages bool
	Name                   string
	DOB                    time.Time
	Gender                 Gender
	PhoneNumber            string
	Website                string
	Bio                    string
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

type Gender uint

const (
	MALE Gender = iota
	FEMALE
)

type Role string

const (
	USER  Role = "ROLE_USER"
	ADMIN Role = "ROLE_ADMIN"
)

type VerificationRequest struct {
	ID                      uuid.UUID `gorm:"primary_key; unique; type:uuid;"`
	UserID                  uuid.UUID
	Name                    string
	Surname                 string
	OfficialDocumentPicture string
	Verified                bool
	VerificationCategory    VerificationCategory
}

type VerificationCategory int

const (
	NEWS_MEDIA VerificationCategory = iota
	SPORTS
	GOVERNMENT_POLITICS
	MUSIC
	FASHION
	ENTERTAINMENT
	BLOGGER_INFLUENCER
	BUSINESS_BRAND_ORGANIZATION
	OTHER
)

type Follow struct {
	UserID      uuid.UUID `gorm:"primary_key; type:uuid;"` //User being followed
	User        User
	FollowerID  uuid.UUID `gorm:"primary_key; type:uuid;"` //Follower of User
	Follower    User
	Muted       bool
	CloseFriend bool `json:"close_friend"`
}

type FollowRequest struct {
	UserID     uuid.UUID `gorm:"primary_key; type:uuid;"` //User requesting a follow
	User       User
	FollowedID uuid.UUID `gorm:"primary_key; type:uuid;"` //User being followed
}

type Block struct {
	UserID    uuid.UUID `gorm:"primary_key; type:uuid;"` //User blocking
	User      User
	BlockedID uuid.UUID `gorm:"primary_key; type:uuid;"` //User being blocked
	Blocked   User
}

type Mute struct {
	UserID  uuid.UUID `gorm:"primary_key; type:uuid;"` //User blocking
	User    User
	MutedID uuid.UUID `gorm:"primary_key; type:uuid;"` //User being blocked
	Muted   User
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}

func (v *VerificationRequest) BeforeCreate(tx *gorm.DB) (err error) {
	v.ID = uuid.New()
	return
}
