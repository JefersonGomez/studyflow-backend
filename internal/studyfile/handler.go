package studyfile

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// UploadStudyFileHandler sube un PDF a un curso
// @Summary      Subir PDF
// @Description  Sube un archivo PDF al curso del usuario autenticado
// @Tags         files
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID del curso"
// @Param        file formData file true "Archivo PDF"
// @Success      201 {object} Studyfile
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /courses/{id}/files [post]
func UploadStudyFileHandler(c *gin.Context) {
	courseID := c.Param("id")
	userID, _ := c.Get("userID")

	// recibir el archivo del request
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no se recibió ningún archivo"})
		return
	}

	// validar que sea PDF
	ext := filepath.Ext(fileHeader.Filename)
	if ext != ".pdf" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "solo se permiten archivos PDF"})
		return
	}

	// crear carpeta del usuario si no existe
	userFolder := fmt.Sprintf("./storage/%s", userID.(string))
	if err := os.MkdirAll(userFolder, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al crear carpeta de almacenamiento"})
		return
	}

	// ruta final donde se guarda el archivo
	storagePath := fmt.Sprintf("%s/%s", userFolder, fileHeader.Filename)

	// guardar el archivo en disco
	if err := c.SaveUploadedFile(fileHeader, storagePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al guardar el archivo"})
		return
	}

	// llamar al service para extraer texto y guardar en BD
	file, err := UploadStudyFileService(userID.(string), courseID, fileHeader.Filename, storagePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, file)
}

// GetCourseFilesHandler obtiene los archivos de un curso
// @Summary      Obtener archivos
// @Description  Obtiene todos los archivos PDF de un curso
// @Tags         files
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID del curso"
// @Success      200 {object} []Studyfile
// @Failure      500 {object} map[string]string
// @Router       /courses/{id}/files [get]
func GetCourseFilesHandler(c *gin.Context) {
	courseID := c.Param("id")

	files, err := GetCourseFilesService(courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, files)
}

// DeleteStudyFileHandler elimina un archivo
// @Summary      Eliminar archivo
// @Description  Elimina un archivo PDF del curso del usuario autenticado
// @Tags         files
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID del archivo"
// @Success      200 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /files/{id} [delete]
func DeleteStudyFileHandler(c *gin.Context) {
	fileID := c.Param("id")
	userID, _ := c.Get("userID")

	err := DeleteStudyFileService(fileID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "archivo eliminado"})
}
