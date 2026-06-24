package ai

import (
	"net/http"
	"strconv"

	"github.com/JefersonGomez/studyflow-backend/internal/note"
	"github.com/JefersonGomez/studyflow-backend/internal/studyfile"
	"github.com/gin-gonic/gin"
)

//funciones http

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

func GenerateQuestionsHandler(c *gin.Context) {
	noteID := c.Param("id")

	n, err := note.GetNotesByID(noteID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "nota no encontrada"})
		return
	}

	questions, err := GenerateQuestions(n.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"questions": questions})
}

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

	c.JSON(http.StatusOK, gin.H{"analysis": analysis})
}
