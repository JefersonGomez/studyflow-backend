package course

import "time"

type Course struct {
	ID          string    `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	UserID      string    `json:"idUser" gorm:"not null"`
	Name        string    `json:"name" gorm:"not null"  `
	Description string    `json:"description"`
	Color       string    `json:"color"`
	CreatedAt   time.Time `json:"createAt" gorm:"autoCreateTime" `
}
