package models

type Attachment struct {
	ID uint `json:"id"`

	StoragePath string `json:"storage_path"`
	PostID      uint   `json:"post_id"`

	Post *Post `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"post,omitempty"`
}

func (Attachment) TableName() string {
	return "post_attachments"
}
