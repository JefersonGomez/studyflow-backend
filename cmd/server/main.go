package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/JefersonGomez/studyflow-backend/docs"
	"github.com/JefersonGomez/studyflow-backend/internal/auth"
	"github.com/JefersonGomez/studyflow-backend/internal/course"
	"github.com/JefersonGomez/studyflow-backend/internal/event"
	"github.com/JefersonGomez/studyflow-backend/internal/note"
	"github.com/JefersonGomez/studyflow-backend/internal/studyfile"
	"github.com/JefersonGomez/studyflow-backend/internal/task"
	"github.com/JefersonGomez/studyflow-backend/internal/user"
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

	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "StudyFlow API corriendo",
			})
		})

		api.GET("/auth/google", auth.GoogleLogin)
		api.GET("/auth/google/callback", auth.GoogleCallback)

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

		courses.GET("/:id/whiteboards", whiteboard.GetCourseWhiteboardsHandler)

		courses.GET("/:id/events", event.GetCourseEventsHandler)

		courses.POST("/:id/files", studyfile.UploadStudyFileHandler)
		courses.GET("/:id/files", studyfile.GetCourseFilesHandler)
	}

	tasks := api.Group("/tasks")
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

	whiteboards := api.Group("/whiteboards")
	whiteboards.Use(middleware.AuthRequired())
	{
		whiteboards.POST("", whiteboard.CreateWhiteboardHandler)
		whiteboards.PUT("/:id", whiteboard.UpdateWhiteboardHandler)
		whiteboards.DELETE("/:id", whiteboard.DeleteWhiteboardHandler)
	}

	events := api.Group("/events")
	events.Use(middleware.AuthRequired())
	{
		events.POST("", event.CreateEventHandler)
		events.GET("", event.GetUserEventsHandler)
		events.PUT("/:id", event.UpdateEventHandler)
		events.DELETE("/:id", event.DeleteEventHandler)
	}

	files := api.Group("/files")
	files.Use(middleware.AuthRequired())
	{
		files.DELETE("/:id", studyfile.DeleteStudyFileHandler)
	}

	port := os.Getenv("PORT")
	fmt.Printf("Servidor corriendo en http://localhost:%s\n", port)
	fmt.Printf("Swagger en http://localhost:%s/swagger/index.html\n", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Error iniciando el servidor")
	}
}
