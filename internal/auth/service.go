package auth

import "github.com/JefersonGomez/studyflow-backend/internal/user"

/* Busca si el usuario ya existe con FindByGoogleID
Si no existe lo crea con CreateUser
Retorna el usuario */

type GoogleUserInfo struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"picture"`
}

func HandlerGoogleUser(googleUser GoogleUserInfo) (*user.User, error) {
	existeUser, err := FindByGoogleID(googleUser.ID)

	if err == nil {
		return existeUser, nil
	}

	newUser := &user.User{
		GoogleID:  googleUser.ID,
		Email:     googleUser.Email,
		Name:      googleUser.Name,
		AvatarURL: googleUser.AvatarURL,
	}

	if err := CreateUser(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}
