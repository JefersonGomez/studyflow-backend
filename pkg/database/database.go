package database

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// controlador para la base de datos postgress
var DB *sqlx.DB

func Connect() {

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sqlx.Connect("postgres", dsn)

	if err != nil {

		log.Fatal("Error al conectar a la base de datos", err)
	}

	//Limita el número máximo de conexiones abiertas simultáneamente a la base de datos.
	db.SetMaxOpenConns(25)

	//Limita el número máximo de conexiones inactivas (idle) que se mantienen en el pool.
	db.SetMaxIdleConns(5)

	DB = db
	log.Println("Conectado a postgress")
}
