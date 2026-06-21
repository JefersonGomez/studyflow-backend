package course

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateCourseRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

// CreateCourseHandler crea una nueva materia
// @Summary      Crear materia
// @Description  Crea una nueva materia para el usuario autenticado
// @Tags         courses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        course body CreateCourseRequest true "Datos de la materia"
// @Success      201 {object} Course
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /courses [post]
func CreateCourseHandler(c *gin.Context) {

	var req CreateCourseRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	userID, _ := c.Get("userID")
	course, err := CreateCourseService(userID.(string), req.Name, req.Description, req.Color)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, course)

}

// @Summary      Obtener materias
// @Description  Obtiene los cursos del usuario autenticado
// @Tags         courses
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} []Course
// @Failure      500 {object} map[string]string
// @Router       /courses [get]
func GetCoursesHandler(c *gin.Context) {

	userID, _ := c.Get("userID")

	courses, err := GetUserCoursesByService(userID.(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, courses)

}

// UpdateCourseHandler actualiza los datos de una materia o curso
// @Summary      Actualizar materia
// @Description  Actualiza los datos de una materia del usuario autenticado
// @Tags         courses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID de la materia"
// @Param        course body CreateCourseRequest true "Datos actualizados de la materia"
// @Success      200 {object} Course
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /courses/{id} [put]
func UpdateCourseHandler(c *gin.Context) {
	courseID := c.Param("id")

	var req CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	course, err := UpdateCourseService(courseID, userID.(string), req.Name, req.Description, req.Color)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, course)
}

// DeleteCourseHandler borra una materia
// @Summary      Eliminar materia
// @Description  Elimina una materia del usuario autenticado
// @Tags         courses
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID de la materia"
// @Success      200 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /courses/{id} [delete]
func DeleteCourseHandler(c *gin.Context) {
	courseID := c.Param("id")
	userID, _ := c.Get("userID")

	err := DeleteCourseService(courseID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "curso eliminado"})
}
