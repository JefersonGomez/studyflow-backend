package note

import (
	"fmt"
	"net/http"

	"github.com/JefersonGomez/studyflow-backend/internal/course"
	"github.com/JefersonGomez/studyflow-backend/pkg/database"
	"github.com/gin-gonic/gin"
)

type CreateNoteRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
}

type NoteResponse struct {
	ID         string  `json:"id"`
	Title      string  `json:"title" binding:"required"`
	Content    string  `json:"content"`
	CourseID   *string `json:"courseId"`
	CourseName string  `json:"courseName"` // ← Nuevo campo
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

func GetAllNotesHandler(c *gin.Context) {
	userID, _ := c.Get("userID")

	var notes []Note
	result := database.DB.Debug().Where("user_id = ?", userID.(string)).Find(&notes)

	fmt.Println("Error:", result.Error)
	fmt.Println("Rows:", result.RowsAffected)
	fmt.Println("Notes:", notes)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, []NoteResponse{})
		return
	}

	// Opcional: Enriquecer con courseName igual que hicimos en tasks/events
	// Por ahora devolvemos las notas puras
	courseIDs := make(map[string]bool)
	for _, t := range notes {
		if t.CourseID != nil && *t.CourseID != "" {
			courseIDs[*t.CourseID] = true
		}
	}

	// 2. Cargar nombres en UNA sola consulta
	coursesMap := make(map[string]string)
	if len(courseIDs) > 0 {
		ids := make([]string, 0, len(courseIDs))
		for id := range courseIDs {
			ids = append(ids, id)
		}

		var courses []course.Course
		// Usamos .Table("courses") si tienes problemas de pluralización,
		// o simplemente database.DB.Find(&courses, ids) si GORM lo maneja bien
		database.DB.Select("id, name").Where("id IN ?", ids).Find(&courses)

		for _, cr := range courses {
			coursesMap[cr.ID] = cr.Name
		}
	}

	// 3. Mapear a respuesta
	response := make([]NoteResponse, len(notes))
	for i, t := range notes {
		courseName := "Sin Materia"
		if t.CourseID != nil && *t.CourseID != "" {
			if name, ok := coursesMap[*t.CourseID]; ok {
				courseName = name
			}
		}

		response[i] = NoteResponse{
			ID:         t.ID,
			Title:      t.Title,
			Content:    t.Content,
			CourseID:   t.CourseID,
			CourseName: courseName,
		}
	}

	c.JSON(http.StatusOK, response)
}
