package whiteboard

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type CreateWhiteboardRequest struct {
	Title    string         `json:"title" binding:"required"`
	CourseID *string        `json:"courseID"`
	Elements datatypes.JSON `json:"elements"`
}

// CreateWhiteboardHandler crea una nueva pizarra
// @Summary      Crear pizarra
// @Description  Crea una nueva pizarra para el usuario autenticado
// @Tags         whiteboards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        whiteboard body CreateWhiteboardRequest true "Datos de la pizarra"
// @Success      201 {object} Whiteboard
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /whiteboards [post]
func CreateWhiteboardHandler(c *gin.Context) {
	var req CreateWhiteboardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")
	courseID := ""
	if req.CourseID != nil {
		courseID = *req.CourseID
	}

	wb, err := CreateWhiteboardService(userID.(string), courseID, req.Title, req.Elements)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, wb)
}

// GetCourseWhiteboardsHandler obtiene las pizarras de un curso
// @Summary      Obtener pizarras de un curso
// @Description  Obtiene todas las pizarras de un curso específico
// @Tags         whiteboards
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID del curso"
// @Success      200 {object} []Whiteboard
// @Failure      500 {object} map[string]string
// @Router       /courses/{id}/whiteboards [get]
func GetCourseWhiteboardsHandler(c *gin.Context) {
	courseID := c.Param("id")

	whiteboards, err := GetCourseWhiteboardsService(courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, whiteboards)
}

// UpdateWhiteboardHandler actualiza una pizarra
// @Summary      Actualizar pizarra
// @Description  Actualiza el contenido de una pizarra del usuario autenticado
// @Tags         whiteboards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID de la pizarra"
// @Param        whiteboard body CreateWhiteboardRequest true "Datos actualizados de la pizarra"
// @Success      200 {object} Whiteboard
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /whiteboards/{id} [put]
func UpdateWhiteboardHandler(c *gin.Context) {
	whiteboardID := c.Param("id")
	userID, _ := c.Get("userID")

	var req CreateWhiteboardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wb, err := UpdateWhiteboardService(whiteboardID, userID.(string), req.Title, req.Elements)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, wb)
}

// DeleteWhiteboardHandler elimina una pizarra
// @Summary      Eliminar pizarra
// @Description  Elimina una pizarra del usuario autenticado
// @Tags         whiteboards
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID de la pizarra"
// @Success      200 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /whiteboards/{id} [delete]
func DeleteWhiteboardHandler(c *gin.Context) {
	whiteboardID := c.Param("id")
	userID, _ := c.Get("userID")

	err := DeleteWhiteboardService(whiteboardID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "pizarra eliminada"})
}
