package studyplan

import "time"

type StudyPlan struct {
	ID        string    `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	CourseID  string    `json:"courseId" gorm:"not null"`
	Content   string    `json:"content" gorm:"not null"`
	Days      int       `json:"days" gorm:"not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
}
