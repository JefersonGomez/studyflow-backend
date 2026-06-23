package studyfile

import "github.com/JefersonGomez/studyflow-backend/pkg/database"

func CreateStudyFile(file *Studyfile) error {

	result := database.DB.Create(file)

	if result.Error != nil {
		return result.Error
	}
	return nil

}

func GetStudyFilesByCourseID(courseID string) ([]Studyfile, error) {

	var files []Studyfile
	result := database.DB.Where("course_id =?", courseID).Find(&files)
	if result.Error != nil {
		return nil, result.Error
	}
	return files, nil

}

func GetStudyFileByID(fileID string) (*Studyfile, error) {

	var file Studyfile

	result := database.DB.Where("id=?", fileID).First(&file)
	if result.Error != nil {
		return nil, result.Error
	}
	return &file, nil

}

func DeleteStudyFile(fileID string) error {

	result := database.DB.Delete(&Studyfile{}, "id = ?", fileID)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
