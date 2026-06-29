package event

import (
	"net/http"
	"time"

	"github.com/JefersonGomez/studyflow-backend/internal/course"
	"github.com/JefersonGomez/studyflow-backend/pkg/database"
	"github.com/gin-gonic/gin"
)

type CreateEventRequest struct {
	Title       string     `json:"title" binding:"required" example:"Examen Final"`
	Description string     `json:"description" example:"Examen del primer parcial"`
	Type        string     `json:"type" binding:"required" example:"exam"`
	StartDate   time.Time  `json:"startDate" binding:"required"`
	EndDate     *time.Time `json:"endDate"`
	CourseID    *string    `json:"courseID" example:"123456"`
}

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

// CreateEventHandler godoc
//
//	@Summary		Crear evento
//	@Description	Crea un nuevo evento para el usuario autenticado
//	@Tags			Events
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			event	body		CreateEventRequest	true	"Datos del evento"
//	@Success		201		{object}	Event
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/events [post]
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

	event, err := CreateEventService(
		userID.(string),
		courseID,
		req.Title,
		req.Description,
		req.Type,
		req.StartDate,
		req.EndDate,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// GetUserEventsHandler godoc
//
//	@Summary		Obtener eventos del usuario
//	@Description	Obtiene todos los eventos del usuario autenticado
//	@Tags			Events
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}		EventResponse
//	@Failure		500	{object}	map[string]string
//	@Router			/events [get]
func GetUserEventsHandler(c *gin.Context) {
	userID, _ := c.Get("userID")

	events, err := GetUserEventsService(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	courseIDs := make(map[string]bool)
	for _, e := range events {
		if e.CourseID != nil && *e.CourseID != "" {
			courseIDs[*e.CourseID] = true
		}
	}

	coursesMap := make(map[string]string)

	if len(courseIDs) > 0 {
		ids := make([]string, 0, len(courseIDs))

		for id := range courseIDs {
			ids = append(ids, id)
		}

		var courses []course.Course

		database.DB.
			Table("courses").
			Select("id, name").
			Where("id IN ?", ids).
			Find(&courses)

		for _, cr := range courses {
			coursesMap[cr.ID] = cr.Name
		}
	}

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

// GetCourseEventsHandler godoc
//
//	@Summary		Obtener eventos por curso
//	@Description	Obtiene todos los eventos asociados a un curso
//	@Tags			Events
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"ID del curso"
//	@Success		200	{array}		Event
//	@Failure		500	{object}	map[string]string
//	@Router			/events/course/{id} [get]
func GetCourseEventsHandler(c *gin.Context) {
	courseID := c.Param("id")

	events, err := GetCourseEventsService(courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

// UpdateEventHandler godoc
//
//	@Summary		Actualizar evento
//	@Description	Actualiza un evento existente
//	@Tags			Events
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string				true	"ID del evento"
//	@Param			event	body		CreateEventRequest	true	"Datos actualizados"
//	@Success		200		{object}	Event
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/events/{id} [put]
func UpdateEventHandler(c *gin.Context) {
	eventID := c.Param("id")
	userID, _ := c.Get("userID")

	var req CreateEventRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event, err := UpdateEventService(
		eventID,
		userID.(string),
		req.Title,
		req.Description,
		req.Type,
		req.StartDate,
		req.EndDate,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, event)
}

// DeleteEventHandler godoc
//
//	@Summary		Eliminar evento
//	@Description	Elimina un evento del usuario autenticado
//	@Tags			Events
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"ID del evento"
//	@Success		200	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/events/{id} [delete]
func DeleteEventHandler(c *gin.Context) {
	eventID := c.Param("id")
	userID, _ := c.Get("userID")

	err := DeleteEventService(eventID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "evento eliminado",
	})
}
