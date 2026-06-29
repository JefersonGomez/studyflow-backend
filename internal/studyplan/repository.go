package studyplan

import "github.com/JefersonGomez/studyflow-backend/pkg/database"

func SaveStudyPlan(courseID string, content string, days int) (*StudyPlan, error) {
	// Reemplazar el plan anterior del mismo curso
	database.DB.Where("course_id = ?", courseID).Delete(&StudyPlan{})

	plan := &StudyPlan{
		CourseID: courseID,
		Content:  content,
		Days:     days,
	}
	result := database.DB.Create(plan)
	return plan, result.Error
}

func GetStudyPlanByCourse(courseID string) (*StudyPlan, error) {
	var plan StudyPlan
	result := database.DB.Where("course_id = ?", courseID).First(&plan)
	return &plan, result.Error
}
