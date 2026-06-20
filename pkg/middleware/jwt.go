package middleware

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string) (string, error) {

	// parsear el token

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		return []byte(os.Getenv("JWT_SECRET")), nil

	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return "", errors.New("no se pudieron leer los datos")
	}

	userID, ok := claims["id"].(string)

	if !ok {

		return "", errors.New("erro en los claims")

	}
	return userID, nil
}
