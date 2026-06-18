package auth

import (
	"github.com/JefersonGomez/studyflow-backend/internal/user"

	"github.com/JefersonGomez/studyflow-backend/pkg/database"
)

func FindByGoogleID(GoogleID string) (*user.User, error) {

	var usuario user.User
	resultado := database.DB.Where("google_id = ?", GoogleID).First(&usuario)

	if resultado.Error != nil {
		return nil, resultado.Error
	}

	return &usuario, nil
}

func CreateUser(user *user.User) error {

	result := database.DB.Create(&user)

	if result.Error != nil {
		return result.Error
	}
	return nil

}
