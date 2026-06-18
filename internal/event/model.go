package event

import "time"

type Event struct {
	ID          string     `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	UserID      string     `json:"userID" gorm:"not null"`
	CourseID    *string    `json:"courseID"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Type        string     `json:"type" gorm:"not null;check:type IN ('examen','quiz','proyecto','laboratorio','clase')"`
	StartDate   time.Time  `json:"startDate" gorm:"not null"`
	EndDate     *time.Time `json:"endDate"`
	CreatedAt   time.Time  `json:"createAt" gorm:"autoCreateTime" `
}
