package task

import (
	"net/http"
	"time"

	"github.com/JefersonGomez/studyflow-backend/internal/course"
	"github.com/JefersonGomez/studyflow-backend/pkg/database"
	"github.com/gin-gonic/gin"
)

type CreateTaskRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	DueDate     time.Time `json:"dueDate" binding:"required"`
}

type TaskResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	DueDate     time.Time `json:"dueDate"`
	CourseID    *string   `json:"courseId"`
	CourseName  string    `json:"courseName"` // ← Nuevo campo
}

// CreateTaskHandler crea una nueva tarea en un curso
// @Summary      Crear tarea
// @Description  Crea una nueva tarea para el curso del usuario autenticado
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID del curso"
// @Param        task body CreateTaskRequest true "Datos de la tarea"
// @Success      201 {object} Task
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /courses/{id}/tasks [post]
func CreateTaskHandler(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	courseID := c.Param("id")
	userID, _ := c.Get("userID")

	task, err := CreateTaskService(userID.(string), courseID, req.Title, req.Description, req.Status, req.DueDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// GetTaskHandler obtiene las tareas de un curso
// @Summary      Obtener tareas
// @Description  Obtiene las tareas de un curso del usuario autenticado
// @Tags         tasks
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID del curso"
// @Success      200 {object} []Task
// @Failure      500 {object} map[string]string
// @Router       /courses/{id}/tasks [get]
func GetTaskHandler(c *gin.Context) {
	courseID := c.Param("id")

	tasks, err := GetCourseTask(courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// UpdateTaskHandler actualiza una tarea
// @Summary      Actualizar tarea
// @Description  Actualiza los datos de una tarea del usuario autenticado
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID de la tarea"
// @Param        task body CreateTaskRequest true "Datos actualizados de la tarea"
// @Success      200 {object} Task
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /tasks/{id} [put]
func UpdateTaskHandler(c *gin.Context) {
	taskID := c.Param("id")

	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	task, err := UpdateTaskService(taskID, userID.(string), "", req.Title, req.Description, req.Status, req.DueDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// DeleteTaskHandler elimina una tarea
// @Summary      Eliminar tarea
// @Description  Elimina una tarea del usuario autenticado
// @Tags         tasks
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID de la tarea"
// @Success      200 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /tasks/{id} [delete]
func DeleteTaskHandler(c *gin.Context) {
	taskID := c.Param("id")
	userID, _ := c.Get("userID")

	err := DeleteTaskService(taskID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "tarea eliminada"})
}

func GetAllTaskHandler(c *gin.Context) {
	userID, _ := c.Get("userID")

	var tasks []Task // Tu modelo original
	database.DB.Where("user_id = ?", userID).Find(&tasks)

	// 1. Recolectar IDs únicos de cursos
	courseIDs := make(map[string]bool)
	for _, t := range tasks {
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
	response := make([]TaskResponse, len(tasks))
	for i, t := range tasks {
		courseName := "Sin Materia"
		if t.CourseID != nil && *t.CourseID != "" {
			if name, ok := coursesMap[*t.CourseID]; ok {
				courseName = name
			}
		}

		response[i] = TaskResponse{
			ID:          t.ID,
			Title:       t.Title,
			Description: t.Description,
			Status:      t.Status,
			DueDate:     t.DueDate,
			CourseID:    t.CourseID,
			CourseName:  courseName,
		}
	}

	c.JSON(http.StatusOK, response)
}
