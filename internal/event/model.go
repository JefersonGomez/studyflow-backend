package event

import "time"

type Event struct {
	ID          string     `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	UserID      string     `json:"userID" gorm:"column:user_id;not null"`
	CourseID    *string    `json:"courseID" gorm:"column:course_id"` // ← Agrega column:course_id
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Type        string     `json:"type" gorm:"column:type;not null;check:type IN ('examen','quiz','proyecto','laboratorio','clase')"`
	StartDate   time.Time  `json:"startDate" gorm:"column:start_date;not null"`      // ← Agrega column:start_date
	EndDate     *time.Time `json:"endDate" gorm:"column:end_date"`                   // ← Agrega column:end_date
	CreatedAt   time.Time  `json:"createAt" gorm:"column:created_at;autoCreateTime"` // ← Agrega column:created_at
}
