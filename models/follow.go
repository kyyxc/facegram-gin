package models

import (
	"time"
)

type Follow struct {
	ID uint `json:"id"`

	FollowerID  uint `json:"follower_id"`
	FollowingID uint `json:"following_id"`
	IsAccepted  bool `json:"is_accepted"`

	CreatedAt time.Time `json:"created_at"`
	
	Follower  *User `gorm:"constraint,OnUpdate:CASCADE,OnDelete:SET NULL;" json:"follower,omitempty"`
	Following *User `gorm:"constraint,OnUpdate:CASCADE,OnDelete:SET NULL" json:"following,omitempty"`
}

func (Follow) TableName() string {
	return "follow"
}
