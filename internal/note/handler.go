package note

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateNoteRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
}

// CreateNoteHandler crea una nueva nota en un curso
// @Summary      Crear nota
// @Description  Crea una nueva nota para el curso del usuario autenticado
// @Tags         notes
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID del curso"
// @Param        note body CreateNoteRequest true "Datos de la nota"
// @Success      201 {object} Note
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /courses/{id}/notes [post]
func CreateNoteHandler(c *gin.Context) {
	var req CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	courseID := c.Param("id")
	userID, _ := c.Get("userID")

	newNote, err := CreateNoteService(userID.(string), courseID, req.Title, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newNote)
}

// GetNotesHandler obtiene las notas de un curso
// @Summary      Obtener notas
// @Description  Obtiene las notas de un curso del usuario autenticado
// @Tags         notes
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID del curso"
// @Success      200 {object} []Note
// @Failure      500 {object} map[string]string
// @Router       /courses/{id}/notes [get]
func GetNotesHandler(c *gin.Context) {
	courseID := c.Param("id")

	notes, err := GetNotesCourseService(courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notes)
}

// UpdateNoteHandler actualiza una nota
// @Summary      Actualizar nota
// @Description  Actualiza los datos de una nota del usuario autenticado
// @Tags         notes
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID de la nota"
// @Param        note body CreateNoteRequest true "Datos actualizados de la nota"
// @Success      200 {object} Note
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /notes/{id} [put]
func UpdateNoteHandler(c *gin.Context) {
	noteID := c.Param("id")
	userID, _ := c.Get("userID")

	var req CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note, err := UpdateNotesService(noteID, userID.(string), "", req.Title, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, note)
}

// DeleteNoteHandler elimina una nota
// @Summary      Eliminar nota
// @Description  Elimina una nota del usuario autenticado
// @Tags         notes
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID de la nota"
// @Success      200 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /notes/{id} [delete]
func DeleteNoteHandler(c *gin.Context) {
	noteID := c.Param("id")
	userID, _ := c.Get("userID")

	err := DeleteNoteService(noteID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "nota eliminada"})
}
