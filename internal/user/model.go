package user

import "time"

type User struct {
	ID        string    `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	GoogleID  string    `json:"google_id" gorm:"unique;not null"`
	Name      string    `json:"name" gorm:"not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	AvatarURL string    `json:"avatarUrl"`
	CreatedAt time.Time `json:"createAt" gorm:"autoCreateTime" `
}
