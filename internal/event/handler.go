package event

import (
	"net/http"
	"time"

	"github.com/JefersonGomez/studyflow-backend/internal/course"
	"github.com/JefersonGomez/studyflow-backend/pkg/database" // Asegúrate de importar tu paquete de DB
	"github.com/gin-gonic/gin"
)

type CreateEventRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Type        string     `json:"type" binding:"required"`
	StartDate   time.Time  `json:"startDate" binding:"required"`
	EndDate     *time.Time `json:"endDate"`
	CourseID    *string    `json:"courseID"`
}

// EventResponse es SOLO para respuestas GET, no modifica el modelo original
type EventResponse struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Type        string     `json:"type"`
	StartDate   time.Time  `json:"startDate"`
	EndDate     *time.Time `json:"endDate"`
	Description string     `json:"description"`
	CourseID    *string    `json:"courseId"`
	CourseName  string     `json:"courseName"`
}

// CreateEventHandler crea un nuevo evento
func CreateEventHandler(c *gin.Context) {
	var req CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	courseID := ""
	if req.CourseID != nil {
		courseID = *req.CourseID
	}

	event, err := CreateEventService(userID.(string), courseID, req.Title, req.Description, req.Type, req.StartDate, req.EndDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// GetUserEventsHandler obtiene todos los eventos del usuario con nombre de curso
func GetUserEventsHandler(c *gin.Context) {
	userID, _ := c.Get("userID")

	events, err := GetUserEventsService(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 1. Recolectar IDs únicos de cursos
	courseIDs := make(map[string]bool)
	for _, e := range events {
		if e.CourseID != nil && *e.CourseID != "" {
			courseIDs[*e.CourseID] = true
		}
	}

	// 2. Cargar nombres usando .Table("courses") para evitar pluralización
	coursesMap := make(map[string]string)
	if len(courseIDs) > 0 {
		ids := make([]string, 0, len(courseIDs))
		for id := range courseIDs {
			ids = append(ids, id)
		}

		// ✅ Usamos el modelo Course real y forzamos la tabla correcta
		var courses []course.Course
		database.DB.Table("courses").Select("id, name").Where("id IN ?", ids).Find(&courses)

		for _, cr := range courses {
			coursesMap[cr.ID] = cr.Name
		}
	}

	// 3. Mapear a EventResponse
	response := make([]EventResponse, len(events))
	for i, e := range events {
		courseName := "Sin Materia"
		if e.CourseID != nil && *e.CourseID != "" {
			if name, ok := coursesMap[*e.CourseID]; ok {
				courseName = name
			}
		}

		response[i] = EventResponse{
			ID:          e.ID,
			Title:       e.Title,
			Type:        e.Type,
			StartDate:   e.StartDate,
			EndDate:     e.EndDate,
			Description: e.Description,
			CourseID:    e.CourseID,
			CourseName:  courseName,
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetCourseEventsHandler obtiene los eventos de un curso específico
func GetCourseEventsHandler(c *gin.Context) {
	courseID := c.Param("id")

	events, err := GetCourseEventsService(courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

// UpdateEventHandler actualiza un evento
func UpdateEventHandler(c *gin.Context) {
	eventID := c.Param("id")
	userID, _ := c.Get("userID")

	var req CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event, err := UpdateEventService(eventID, userID.(string), req.Title, req.Description, req.Type, req.StartDate, req.EndDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, event)
}

// DeleteEventHandler elimina un evento
func DeleteEventHandler(c *gin.Context) {
	eventID := c.Param("id")
	userID, _ := c.Get("userID")

	err := DeleteEventService(eventID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "evento eliminado"})
}
