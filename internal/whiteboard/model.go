package whiteboard

import (
	"time"

	"gorm.io/datatypes"
)

type Whiteboard struct {
	ID        string         `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	UserID    string         `json:"userID" gorm:"not null"`
	CourseID  *string        `json:"courseID"`
	Title     string         `json:"title"`
	Elements  datatypes.JSON `json:"elemts"`
	CreatedAt time.Time      `json:"createAt" gorm:"autoCreateTime" `
	UptadeAt  time.Time      `json:"updateAt" gorm:"autoCreateTime"`
}
