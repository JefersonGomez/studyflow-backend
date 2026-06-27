package note

import "time"

type Note struct {
	ID        string    `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	UserID    string    `json:"userID" gorm:"not null"`
	CourseID  *string   `json:"courseID" gorm:"not null"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createAt" gorm:"autoCreateTime" `
	UpdatedAt time.Time `json:"updateAt" gorm:"autoUpdateTime;column:updated_at"`
}
