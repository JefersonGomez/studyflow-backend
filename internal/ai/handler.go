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
	"github.com/JefersonGomez/studyflow-backend/internal/studyplan"
	"github.com/JefersonGomez/studyflow-backend/pkg/database"
	"github.com/gin-gonic/gin"
)

// funciones http
// SummarizeNoteHandler godoc
//
//	@Summary		Resumir nota con IA
//	@Description	Genera un resumen automático del contenido de una nota utilizando IA
//	@Tags			AI
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"ID de la nota"
//	@Success		200	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/ai/notes/{id}/summary [get]
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

// GenerateStudyPlanHandler godoc
//
//	@Summary		Generar plan de estudio desde archivo
//	@Description	Genera un plan de estudio basado en el contenido de un PDF o material de estudio
//	@Tags			AI
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string	true	"ID del archivo"
//	@Param			days	query		int		true	"Cantidad de días"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/ai/files/{id}/study-plan [get]
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

// GenerateQuestionsHandler godoc
//
//	@Summary		Generar preguntas de estudio
//	@Description	Genera preguntas de práctica a partir de una nota usando IA
//	@Tags			AI
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"ID de la nota"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/ai/notes/{id}/questions [get]
func GenerateQuestionsHandler(c *gin.Context) {
	noteID := c.Param("id")

	n, err := note.GetNotesByID(noteID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "nota no encontrada"})
		return
	}

	raw, err := GenerateQuestions(n.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Parsear el string JSON que devuelve la IA
	var questions []map[string]string
	if err := json.Unmarshal([]byte(raw), &questions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "respuesta de IA no válida"})
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

// AnalyzePDFHandler godoc
//
//	@Summary		Analizar PDF con IA
//	@Description	Analiza un PDF, extrae fechas importantes y crea eventos automáticamente en el calendario
//	@Tags			AI
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"ID del archivo"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/ai/files/{id}/analyze [post]
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

// GenerateStudyPlanByCourseHandler godoc
//
//	@Summary		Generar plan de estudio por curso
//	@Description	Genera y guarda un plan de estudio usando todos los archivos asociados a un curso
//	@Tags			AI
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string	true	"ID del curso"
//	@Param			days	query		int		true	"Cantidad de días"
//	@Success		200		{object}	studyplan.StudyPlan
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/ai/courses/{id}/study-plan [post]
func GenerateStudyPlanByCourseHandler(c *gin.Context) {
	courseID := c.Param("id")

	files, err := studyfile.GetStudyFilesByCourseID(courseID)
	if err != nil || len(files) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no hay archivos en este curso"})
		return
	}

	daysStr := c.Query("days")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "days debe ser un número mayor a 0"})
		return
	}

	combinedText := ""
	for _, f := range files {
		combinedText += f.ParsedText + "\n\n"
	}

	content, err := GenerateStudyPlan(combinedText, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ✅ Guardar en BD
	plan, err := studyplan.SaveStudyPlan(courseID, content, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al guardar el plan"})
		return
	}

	c.JSON(http.StatusOK, plan)
}

// GetStudyPlanHandler godoc
//
//	@Summary		Obtener plan de estudio
//	@Description	Obtiene el último plan de estudio generado para un curso
//	@Tags			AI
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"ID del curso"
//	@Success		200	{object}	studyplan.StudyPlan
//	@Failure		404	{object}	map[string]string
//	@Router			/ai/courses/{id}/study-plan [get]
func GetStudyPlanHandler(c *gin.Context) {
	courseID := c.Param("id")

	plan, err := studyplan.GetStudyPlanByCourse(courseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no hay plan generado para este curso"})
		return
	}

	c.JSON(http.StatusOK, plan)
}
