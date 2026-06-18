package database

import (
	"fmt"
	"log"
	"os"

	"github.com/JefersonGomez/studyflow-backend/internal/course"
	"github.com/JefersonGomez/studyflow-backend/internal/event"
	"github.com/JefersonGomez/studyflow-backend/internal/note"
	"github.com/JefersonGomez/studyflow-backend/internal/studyfile"
	"github.com/JefersonGomez/studyflow-backend/internal/task"
	"github.com/JefersonGomez/studyflow-backend/internal/user"
	"github.com/JefersonGomez/studyflow-backend/internal/whiteboard"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// controlador para la base de datos postgress
var DB *gorm.DB

func Connect() {

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {

		log.Fatal("Error al conectar a la base de datos", err)
	}
	DB = db
	log.Println("Conectado a postgress")
}

func Migrate() {
	DB.AutoMigrate(
		&user.User{},
		&course.Course{},
		&event.Event{},
		&task.Task{},
		&note.Note{},
		&whiteboard.Whiteboard{},
		&studyfile.Studyfile{},
	)

}
