package auth

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

/* 1. Frontend llama a /auth/google
2. Backend redirige a Google
3. Google redirige a /auth/google/callback con un código
4. Backend intercambia el código por un token
5. Backend obtiene los datos del usuario de Google
6. Backend crea o busca el usuario en la BD
7. Backend genera un JWT y lo devuelve al frontend */

var googleOauthConfig = &oauth2.Config{

	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	},
	Endpoint: google.Endpoint,
}

// redirige a google(redirige al usuario a la pantalla de login de google)
func GoogleLogin(c *gin.Context) {
	//AuthCode genera una url de goole
	url := googleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)

}

// maneja la respuesta de goole()
func GoogleCallback(c *gin.Context) {
	code := c.Query("code")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Codigo no recibido"})
		return
	}
	//exchage intercabia el codigo por un token de acceso real
	token, err := googleOauthConfig.Exchange(c, code)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al intercambiar el código"})

	}
	_ = token // por ahora, aquí llamaremos al service

}
