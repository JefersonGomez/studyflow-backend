package ai

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/JefersonGomez/studyflow-backend/internal/event"
	"github.com/JefersonGomez/studyflow-backend/internal/note"
	"github.com/JefersonGomez/studyflow-backend/internal/studyfile"
	"github.com/JefersonGomez/studyflow-backend/pkg/database"
	"github.com/gin-gonic/gin"
)

//funciones http

func SummarizeNoteHandler(c *gin.Context) {
	noteID := c.Param("id")

	n, err := note.GetNotesByID(noteID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "nota no encontrada"})
		return
	}

	summary, err := SummarizeNote(n.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"summary": summary})
}

func GenerateStudyPlanHandler(c *gin.Context) {
	fileID := c.Param("id")

	file, err := studyfile.GetStudyFileByID(fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "archivo no encontrado"})
		return
	}

	daysStr := c.Query("days")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "el parámetro days debe ser un número mayor a 0"})
		return
	}

	studyPlan, err := GenerateStudyPlan(file.ParsedText, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"studyPlan": studyPlan})
}

func GenerateQuestionsHandler(c *gin.Context) {
	noteID := c.Param("id")

	n, err := note.GetNotesByID(noteID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "nota no encontrada"})
		return
	}

	questions, err := GenerateQuestions(n.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"questions": questions})
}

func mapEventType(tipo string) string {
	switch tipo {
	case "examen":
		return "examen"
	case "quiz":
		return "quiz"
	case "proyecto", "propuesta", "avance", "seguimiento":
		return "proyecto"
	case "tarea", "investigación", "investigacion", "laboratorio":
		return "laboratorio"
	default:
		return "clase"
	}
}
func AnalyzePDFHandler(c *gin.Context) {
	fileID := c.Param("id")

	file, err := studyfile.GetStudyFileByID(fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "archivo no encontrado"})
		return
	}

	analysis, err := AnalyzePDF(file.ParsedText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Limpiar el JSON de posibles backticks que la IA agrega
	clean := strings.TrimSpace(analysis)
	clean = strings.TrimPrefix(clean, "```json")
	clean = strings.TrimPrefix(clean, "```")
	clean = strings.TrimSuffix(clean, "```")
	clean = strings.TrimSpace(clean)

	// Parsear el JSON de eventos
	var result struct {
		Eventos []struct {
			Titulo string `json:"titulo"`
			Fecha  string `json:"fecha"`
			Tipo   string `json:"tipo"`
		} `json:"eventos"`
	}

	if err := json.Unmarshal([]byte(clean), &result); err != nil {
		// Si no se puede parsear, devolver el texto crudo
		c.JSON(http.StatusOK, gin.H{"analysis": analysis, "eventos": []interface{}{}})
		return
	}

	// Crear los eventos en la base de datos
	userID, _ := c.Get("userID")
	courseID := file.CourseID
	createdCount := 0

	for _, e := range result.Eventos {
		fecha, err := time.Parse("2006-01-02", e.Fecha)
		if err != nil {
			continue
		}

		newEvent := event.Event{
			UserID:      userID.(string),
			CourseID:    courseID,
			Title:       e.Titulo,
			Type:        mapEventType(e.Tipo),
			StartDate:   fecha,
			Description: "Generado automáticamente por IA",
		}

		database.DB.Create(&newEvent)
		createdCount++
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Se crearon %d eventos en el calendario", createdCount),
		"eventos": result.Eventos,
		"creados": createdCount,
	})
}
