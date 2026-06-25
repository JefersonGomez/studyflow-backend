package event

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateEventRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Type        string     `json:"type" binding:"required"`
	StartDate   time.Time  `json:"startDate" binding:"required"`
	EndDate     *time.Time `json:"endDate"`
	CourseID    *string    `json:"courseID"` // ← Agrega este campo
}

// CreateEventHandler crea un nuevo evento
// @Summary      Crear evento
// @Description  Crea un nuevo evento para el usuario autenticado
// @Tags         events
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        event body CreateEventRequest true "Datos del evento"
// @Success      201 {object} Event
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /events [post]
func CreateEventHandler(c *gin.Context) {
	var req CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	// ← CAMBIO AQUÍ: Obtener courseID del body, no de query params
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

// GetUserEventsHandler obtiene todos los eventos del usuario
// @Summary      Obtener eventos del usuario
// @Description  Obtiene todos los eventos del usuario autenticado
// @Tags         events
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} []Event
// @Failure      500 {object} map[string]string
// @Router       /events [get]
func GetUserEventsHandler(c *gin.Context) {
	userID, _ := c.Get("userID")

	events, err := GetUserEventsService(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

// GetCourseEventsHandler obtiene los eventos de un curso
// @Summary      Obtener eventos de un curso
// @Description  Obtiene todos los eventos de un curso específico
// @Tags         events
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID del curso"
// @Success      200 {object} []Event
// @Failure      500 {object} map[string]string
// @Router       /courses/{id}/events [get]
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
// @Summary      Actualizar evento
// @Description  Actualiza los datos de un evento del usuario autenticado
// @Tags         events
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID del evento"
// @Param        event body CreateEventRequest true "Datos actualizados del evento"
// @Success      200 {object} Event
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /events/{id} [put]
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
// @Summary      Eliminar evento
// @Description  Elimina un evento del usuario autenticado
// @Tags         events
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID del evento"
// @Success      200 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /events/{id} [delete]
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
