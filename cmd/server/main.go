package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"

	_ "github.com/JefersonGomez/studyflow-backend/docs"
	"github.com/JefersonGomez/studyflow-backend/internal/auth"
	"github.com/JefersonGomez/studyflow-backend/internal/course"
	"github.com/JefersonGomez/studyflow-backend/internal/event"
	"github.com/JefersonGomez/studyflow-backend/internal/note"
	"github.com/JefersonGomez/studyflow-backend/internal/studyfile"
	"github.com/JefersonGomez/studyflow-backend/internal/task"
	"github.com/JefersonGomez/studyflow-backend/internal/whiteboard"
	"github.com/JefersonGomez/studyflow-backend/pkg/database"
	"github.com/JefersonGomez/studyflow-backend/pkg/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title studyFlow API
// @version 1.0
// @description API backend para la plataforma StudyFlow
// @host localhost:3000
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {

	// CARGAR VARIABLES DE ENTORNO

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error al cargar el .env")

	}

	database.Connect()

	database.DB.AutoMigrate(
		&user.User{},
		&course.Course{},
		&event.Event{},
		&task.Task{},
		&note.Note{},
		&whiteboard.Whiteboard{},
		&studyfile.Studyfile{},
	)

	// Modo de gin segun entorno
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	//CORS : permite request desde el frontend
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// swager

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// rutas base
	api := r.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {

			c.JSON(200, gin.H{
				"status":  "ok",
				"massage": "Study API corriendo",
			})

		})

		api.GET("/auth/google", auth.GoogleLogin)
		api.GET("/auth/google/callback", auth.GoogleCallback)

		// ejemplo de ruta protegida por un middleware
		api.GET("/me", middleware.AuthRequired(), func(c *gin.Context) {
			userID, _ := c.Get("userID")
			c.JSON(http.StatusOK, gin.H{"userID": userID})
		})

	}
	courses := api.Group("/courses")
	courses.Use(middleware.AuthRequired())
	{
		courses.POST("", course.CreateCourseHandler)
		courses.GET("", course.GetCoursesHandler)
		courses.PUT("/:id", course.UpdateCourseHandler)
		courses.DELETE("/:id", course.DeleteCourseHandler)
		courses.POST("/:id/tasks", task.CreateTaskHandler)
		courses.GET("/:id/tasks", task.GetTaskHandler)
		courses.POST("/:id/notes", note.CreateNoteHandler)
		courses.GET("/:id/notes", note.GetNotesHandler)
	}

	tasks := api.Group("/task")

	tasks.Use(middleware.AuthRequired())
	{

		tasks.PUT("/:id", task.UpdateTaskHandler)
		tasks.DELETE("/:id", task.DeleteTaskHandler)
	}

	notes := api.Group("/notes")
	notes.Use(middleware.AuthRequired())
	{
		notes.PUT("/:id", note.UpdateNoteHandler)
		notes.DELETE("/:id", note.DeleteNoteHandler)
	}

	port := os.Getenv("PORT")

	fmt.Printf("Servidor corriendo en http://localhost:%s\n", port)
	fmt.Printf("Swagger en http://localhost:%s/swagger/index.html\n", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Error iniciando el servidor")
	}

}
