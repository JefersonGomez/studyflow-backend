package auth

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/JefersonGomez/studyflow-backend/pkg/middleware"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// getGoogleOauthConfig construye la config en el momento que se llama,
// no cuando el paquete se carga. Así garantizamos que el .env ya esté cargado.
func getGoogleOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func GoogleLogin(c *gin.Context) {
	config := getGoogleOauthConfig()
	url := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c *gin.Context) {
	config := getGoogleOauthConfig()
	code := c.Query("code")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Codigo no recibido"})
		return
	}

	token, err := config.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al intercambiar el código"})
		return
	}

	client := config.Client(c, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener el token del usuario"})
		return
	}
	defer resp.Body.Close()

	var googleUser GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al leer datos del usuario"})
		return
	}

	u, err := HandlerGoogleUser(googleUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al guardar los datos del usuario"})
		return
	}

	jwtToken, err := middleware.GenerateJWT(u.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"eror": "no se puedo generar el token"})
		return

	}

	c.Redirect(http.StatusTemporaryRedirect, "http://localhost:5173/auth/callback?token="+jwtToken)
}
