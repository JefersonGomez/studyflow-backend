package course

import (
	"errors"
	"fmt"

	"github.com/JefersonGomez/studyflow-backend/pkg/database"
	"gorm.io/gorm"
)

func CreateCourse(course *Course) error {

	result := database.DB.Create(course)

	if result.Error != nil {
		return result.Error
	}
	return nil

}

func GetCoursesByUserID(userID string) ([]Course, error) {

	var courses []Course
	result := database.DB.Where("user_id = ?", userID).Find(&courses)

	if result.Error != nil {
		return nil, result.Error
	}
	return courses, nil

}
func GetCourseByID(courseID string) (*Course, error) {

	var course Course

	result := database.DB.Where("id = ?", courseID).First(&course)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return &Course{}, fmt.Errorf("el curso con id %s no existe", courseID)

		}
		return &Course{}, result.Error
	}

	return &course, nil
}
func UpdateCourse(course *Course) error {

	result := database.DB.Save(course)
	if result.Error != nil {
		return result.Error
	}
	return nil

}
func DeleteCourse(courseID string) error {

	result := database.DB.Delete(&Course{}, "id=?", courseID)

	if result.Error != nil {

		return result.Error
	}
	return nil

}
