package studyfile

import "time"

type Studyfile struct {
	ID          string    `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	UserID      string    `json:"userID" gorm:"not null"`
	CourseID    *string   `json:"courseID"`
	FileName    string    `json:"fileName"`
	StoragePath string    `json:"storageFile"`
	ParsedText  string    `json:"parseText"`
	AiProcessed bool      `json:"iaProcessed"`
	CreatedAt   time.Time `json:"createAt" gorm:"autoCreateTime" `
}
