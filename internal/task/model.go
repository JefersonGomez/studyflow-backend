package task

import "time"

type Task struct {
	ID          string    `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	UserID      string    `json:"userID" gorm:"not null"`
	CourseID    *string   `json:"courseID"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	DueDate     time.Time `json:"dueDate" gorm:"not null"`
	CreatedAt   time.Time `json:"createAt" gorm:"autoCreateTime" `
	UptadeAt    time.Time `json:"updateAt" gorm:"autoUpdateTime"`
}
